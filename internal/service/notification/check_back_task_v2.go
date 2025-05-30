package notification

import (
	"context"
	"fmt"
	"time"

	clientv1 "gitee.com/flycash/notification-platform/api/proto/gen/client/v1"
	"gitee.com/flycash/notification-platform/internal/domain"
	"gitee.com/flycash/notification-platform/internal/pkg/grpc"
	"gitee.com/flycash/notification-platform/internal/pkg/loopjob"
	"gitee.com/flycash/notification-platform/internal/pkg/sharding"
	"gitee.com/flycash/notification-platform/internal/repository"
	"gitee.com/flycash/notification-platform/internal/service/config"
	"github.com/ecodeclub/ekit/list"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gotomicro/ego/client/egrpc"
	"github.com/hashicorp/go-multierror"
	"github.com/meoying/dlock-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type TxCheckTaskV2 struct {
	repo      repository.TxNotificationRepository
	txnStr    sharding.ShardingStrategy
	nStr      sharding.ShardingStrategy
	configSvc config.BusinessConfigService
	lock      dlock.Client
	batchSize int
	clients   *grpc.Clients[clientv1.TransactionCheckServiceClient]
	sem       loopjob.ResourceSemaphore
}

func NewTxCheckTaskV2(repo repository.TxNotificationRepository, configSvc config.BusinessConfigService,
	lock dlock.Client, txnStr, nStr sharding.ShardingStrategy,
	sem loopjob.ResourceSemaphore,
) *TxCheckTaskV2 {
	return &TxCheckTaskV2{
		repo:      repo,
		configSvc: configSvc,
		clients: grpc.NewClients[clientv1.TransactionCheckServiceClient](func(conn *egrpc.Component) clientv1.TransactionCheckServiceClient {
			return clientv1.NewTransactionCheckServiceClient(conn)
		}),
		lock:      lock,
		txnStr:    txnStr,
		nStr:      nStr,
		batchSize: defaultBatchSize,
		sem:       sem,
	}
}

func (task *TxCheckTaskV2) Start(ctx context.Context) {
	const key = "notification_check_task_v2"
	go loopjob.NewShardingLoopJob(task.lock, key, task.oneLoop, task.nStr, task.sem).Run(ctx)
}

//nolint:dupl // 为了演示
func (task *TxCheckTaskV2) oneLoop(ctx context.Context) error {
	loopCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	txNotifications, err := task.repo.FindCheckBack(loopCtx, 0, task.batchSize)
	if err != nil {
		return err
	}

	if len(txNotifications) == 0 {
		// 避免立刻又调度
		time.Sleep(time.Second)
		return nil
	}

	bizIDs := slice.Map(txNotifications, func(_ int, src domain.TxNotification) int64 {
		return src.BizID
	})
	configMap, err := task.configSvc.GetByIDs(loopCtx, bizIDs)
	if err != nil {
		return err
	}
	length := len(txNotifications)
	// 这一次回查没拿到明确结果的
	retryTxns := &list.ConcurrentList[domain.TxNotification]{
		List: list.NewArrayList[domain.TxNotification](length),
	}

	// 要回滚的
	failTxns := &list.ConcurrentList[domain.TxNotification]{
		List: list.NewArrayList[domain.TxNotification](length),
	}
	// 要提交的
	commitTxns := &list.ConcurrentList[domain.TxNotification]{
		List: list.NewArrayList[domain.TxNotification](length),
	}
	var eg errgroup.Group
	for idx := range txNotifications {
		eg.Go(func() error {
			// 并发去回查
			txNotification := txNotifications[idx]
			// 我在这里发起了回查，而后拿到了结果
			txn := task.oneBackCheck(loopCtx, configMap, txNotification)
			switch txn.Status {
			case domain.TxNotificationStatusPrepare:
				// 查到还是 Prepare 状态
				_ = retryTxns.Append(txn)
			case domain.TxNotificationStatusFail, domain.TxNotificationStatusCancel:
				_ = failTxns.Append(txn)
			case domain.TxNotificationStatusCommit:
				_ = commitTxns.Append(txn)
			default:
				return errors.New("不合法的回查状态")
			}
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}
	// 挨个处理，更新数据库状态
	// 数据库就可以一次性执行完，规避频繁更新数据库
	err = task.updateStatus(loopCtx, retryTxns, domain.SendStatusPrepare)
	err = multierror.Append(err, task.updateStatus(loopCtx, failTxns, domain.SendStatusFailed))
	// 转 PENDING，后续 Scheduler 会调度执行
	err = multierror.Append(err, task.updateStatus(loopCtx, commitTxns, domain.SendStatusPending))
	return err
}

// 校验完了
func (task *TxCheckTaskV2) oneBackCheck(ctx context.Context, configMap map[int64]domain.BusinessConfig, txNotification domain.TxNotification) domain.TxNotification {
	bizConfig, ok := configMap[txNotification.BizID]
	if !ok || bizConfig.TxnConfig == nil {
		// 没设置，不需要回查
		txNotification.NextCheckTime = 0
		txNotification.Status = domain.TxNotificationStatusFail
		return txNotification
	}

	txConfig := bizConfig.TxnConfig
	// 发起回查
	res, err := task.getCheckBackRes(ctx, *txConfig, txNotification)
	// 执行了一次回查，要 +1
	txNotification.CheckCount++
	// 回查失败了
	if err != nil || res == unknownStatus {
		// 重新计算下一次的回查时间
		txNotification.SetNextCheckBackTimeAndStatus(txConfig)
		return txNotification
	}
	switch res {
	case cancelStatus:
		txNotification.NextCheckTime = 0
		txNotification.Status = domain.TxNotificationStatusCancel
	case committedStatus:
		txNotification.NextCheckTime = 0
		txNotification.Status = domain.TxNotificationStatusCommit
	}
	return txNotification
}

func (task *TxCheckTaskV2) getCheckBackRes(ctx context.Context, conf domain.TxnConfig, txn domain.TxNotification) (status int, err error) {
	defer func() {
		if r := recover(); r != nil {
			if str, ok := r.(string); ok {
				err = errors.New(str)
			} else {
				err = fmt.Errorf("未知panic类型: %v", r)
			}
		}
	}()
	// 借助服务发现来回查
	client := task.clients.Get(conf.ServiceName)

	req := &clientv1.TransactionCheckServiceCheckRequest{Key: txn.Key}
	resp, err := client.Check(ctx, req)
	if err != nil {
		return unknownStatus, err
	}
	return int(resp.Status), nil
}

func (task *TxCheckTaskV2) updateStatus(ctx context.Context,
	list *list.ConcurrentList[domain.TxNotification], status domain.SendStatus,
) error {
	if list.Len() == 0 {
		return nil
	}
	txns := list.AsSlice()
	return task.repo.UpdateCheckStatus(ctx, txns, status)
}
