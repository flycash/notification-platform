package limit

import (
	"context"
	"sync/atomic"

	"gitee.com/flycash/notification-platform/internal/errs"
	"gitee.com/flycash/notification-platform/internal/pkg/config"
	"gitee.com/flycash/notification-platform/internal/pkg/ratelimit"
	"github.com/gotomicro/ego/core/elog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimitFailoverBuilder 限流后触发故障转移
type RateLimitFailoverBuilder struct {
	svcInstance config.ServiceInstance

	limitedKey     string
	limiter        ratelimit.Limiter
	inLimitedState atomic.Bool

	failoverMgr config.FailoverManager

	logger *elog.Component
}

func NewRateLimitFailoverBuilder(
	svcInstance config.ServiceInstance,
	limitedKey string,
	limiter ratelimit.Limiter,
	failoverMgr config.FailoverManager,
) *RateLimitFailoverBuilder {
	return &RateLimitFailoverBuilder{
		svcInstance: svcInstance,
		limitedKey:  limitedKey,
		limiter:     limiter,
		failoverMgr: failoverMgr,
		logger:      elog.DefaultLogger,
	}
}

func (b *RateLimitFailoverBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		limited, err := b.limiter.Limit(ctx, b.limitedKey)
		if err != nil {
			// 保守策略
			return nil, status.Errorf(codes.ResourceExhausted, "%s", errs.ErrRateLimited)
		}
		if limited {
			// 触发限流标记限流状态
			if b.inLimitedState.CompareAndSwap(false, true) {
				// 标记故障转移
				err1 := b.failoverMgr.Failover(ctx, b.svcInstance)
				if err1 != nil {
					// 只记录日志
					b.logger.Error("故障转移失败",
						elog.FieldErr(err1),
						elog.Any("ServiceInstance", b.svcInstance),
					)
				}
			}
			// 请求被限流
			return nil, status.Errorf(codes.ResourceExhausted, "%s", errs.ErrRateLimited)
		}

		// 之前处于限流状态
		if b.inLimitedState.CompareAndSwap(true, false) {
			// 限流解除后，标记已恢复
			err2 := b.failoverMgr.Recover(ctx, b.svcInstance)
			if err2 != nil {
				// 只记录日志
				b.logger.Error("故障恢复失败",
					elog.FieldErr(err2),
					elog.Any("ServiceInstance", b.svcInstance),
				)
			}
		}

		return handler(ctx, req)
	}
}
