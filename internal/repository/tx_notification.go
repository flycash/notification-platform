package repository

import (
	"context"

	"gitee.com/flycash/notification-platform/internal/domain"
	"gitee.com/flycash/notification-platform/internal/repository/dao"
)

type TxNotificationRepository interface {
	Create(ctx context.Context, notification domain.TxNotification) (uint64, error)
	FindCheckBack(ctx context.Context, offset, limit int) ([]domain.TxNotification, error)
	UpdateStatus(ctx context.Context, bizID int64, key string, status domain.TxNotificationStatus, notificationStatus domain.SendStatus) error
	UpdateCheckStatus(ctx context.Context, txNotifications []domain.TxNotification, notificationStatus domain.SendStatus) error
}

type txNotificationRepo struct {
	txdao dao.TxNotificationDAO
}

// NewTxNotificationRepository creates a new TxNotificationRepository instance
func NewTxNotificationRepository(txdao dao.TxNotificationDAO) TxNotificationRepository {
	return &txNotificationRepo{
		txdao: txdao,
	}
}

func (t *txNotificationRepo) First(ctx context.Context, txID int64) (domain.TxNotification, error) {
	noti, err := t.txdao.First(ctx, txID)
	if err != nil {
		return domain.TxNotification{}, err
	}
	return t.toDomain(noti), nil
}

func (t *txNotificationRepo) Create(ctx context.Context, txn domain.TxNotification) (uint64, error) {
	// 转换领域模型到DAO对象
	txnEntity := t.toDao(txn)
	notificationEntity := t.toEntity(txn.Notification)
	// 调用DAO层创建记录
	return t.txdao.Prepare(ctx, txnEntity, notificationEntity)
}

// toEntity 将领域对象转换为DAO实体
func (t *txNotificationRepo) toEntity(notification domain.Notification) dao.Notification {
	templateParams, _ := notification.MarshalTemplateParams()
	receivers, _ := notification.MarshalReceivers()
	return dao.Notification{
		ID:                notification.ID,
		BizID:             notification.BizID,
		Key:               notification.Key,
		Receivers:         receivers,
		Channel:           string(notification.Channel),
		TemplateID:        notification.Template.ID,
		TemplateVersionID: notification.Template.VersionID,
		TemplateParams:    templateParams,
		Status:            string(notification.Status),
		ScheduledSTime:    notification.ScheduledSTime.UnixMilli(),
		ScheduledETime:    notification.ScheduledETime.UnixMilli(),
		Version:           notification.Version,
	}
}

func (t *txNotificationRepo) FindCheckBack(ctx context.Context, offset, limit int) ([]domain.TxNotification, error) {
	// 调用DAO层查询记录
	daoNotifications, err := t.txdao.FindCheckBack(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// 将DAO对象列表转换为领域模型列表
	result := make([]domain.TxNotification, 0, len(daoNotifications))
	for _, daoNotification := range daoNotifications {
		result = append(result, t.toDomain(daoNotification))
	}
	return result, nil
}

func (t *txNotificationRepo) UpdateStatus(ctx context.Context, bizID int64, key string, status domain.TxNotificationStatus, notificationStatus domain.SendStatus) error {
	// 直接调用DAO层更新状态
	return t.txdao.UpdateStatus(ctx, bizID, key, status, notificationStatus)
}

func (t *txNotificationRepo) UpdateCheckStatus(ctx context.Context, txNotifications []domain.TxNotification, notificationStatus domain.SendStatus) error {
	// 将领域模型列表转换为DAO对象列表
	daoNotifications := make([]dao.TxNotification, 0, len(txNotifications))
	for idx := range txNotifications {
		txNotification := txNotifications[idx]
		daoNotifications = append(daoNotifications, t.toDao(txNotification))
	}

	// 调用DAO层更新检查状态
	return t.txdao.UpdateCheckStatus(ctx, daoNotifications, notificationStatus)
}

// toDomain 将DAO对象转换为领域模型
func (t *txNotificationRepo) toDomain(daoNotification dao.TxNotification) domain.TxNotification {
	return domain.TxNotification{
		TxID: daoNotification.TxID,
		Notification: domain.Notification{
			ID: daoNotification.NotificationID,
		},
		Key:           daoNotification.Key,
		BizID:         daoNotification.BizID,
		Status:        domain.TxNotificationStatus(daoNotification.Status),
		CheckCount:    daoNotification.CheckCount,
		NextCheckTime: daoNotification.NextCheckTime,
		Ctime:         daoNotification.Ctime,
		Utime:         daoNotification.Utime,
	}
}

// toDao 将领域模型转换为DAO对象
func (t *txNotificationRepo) toDao(domainNotification domain.TxNotification) dao.TxNotification {
	return dao.TxNotification{
		TxID:           domainNotification.TxID,
		Key:            domainNotification.Key,
		NotificationID: domainNotification.Notification.ID,
		BizID:          domainNotification.BizID,
		Status:         string(domainNotification.Status),
		CheckCount:     domainNotification.CheckCount,
		NextCheckTime:  domainNotification.NextCheckTime,
		Ctime:          domainNotification.Ctime,
		Utime:          domainNotification.Utime,
	}
}
