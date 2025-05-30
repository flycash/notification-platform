// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"time"

	"gitee.com/flycash/notification-platform/internal/api/grpc"
	"gitee.com/flycash/notification-platform/internal/domain"
	ioc2 "gitee.com/flycash/notification-platform/internal/ioc"
	"gitee.com/flycash/notification-platform/internal/repository"
	"gitee.com/flycash/notification-platform/internal/repository/cache/local"
	"gitee.com/flycash/notification-platform/internal/repository/cache/redis"
	"gitee.com/flycash/notification-platform/internal/repository/dao"
	"gitee.com/flycash/notification-platform/internal/service/audit"
	"gitee.com/flycash/notification-platform/internal/service/channel"
	"gitee.com/flycash/notification-platform/internal/service/config"
	"gitee.com/flycash/notification-platform/internal/service/notification"
	"gitee.com/flycash/notification-platform/internal/service/notification/callback"
	"gitee.com/flycash/notification-platform/internal/service/provider"
	"gitee.com/flycash/notification-platform/internal/service/provider/manage"
	"gitee.com/flycash/notification-platform/internal/service/provider/sequential"
	"gitee.com/flycash/notification-platform/internal/service/provider/sms"
	"gitee.com/flycash/notification-platform/internal/service/provider/sms/client"
	"gitee.com/flycash/notification-platform/internal/service/quota"
	"gitee.com/flycash/notification-platform/internal/service/scheduler"
	"gitee.com/flycash/notification-platform/internal/service/sender"
	"gitee.com/flycash/notification-platform/internal/service/sendstrategy"
	manage2 "gitee.com/flycash/notification-platform/internal/service/template/manage"
	"gitee.com/flycash/notification-platform/internal/test/ioc"
	"github.com/ecodeclub/ekit/pool"
	"github.com/google/wire"
	"github.com/gotomicro/ego/core/econf"
)

// Injectors from wire.go:

func InitGrpcServer(clients map[string]client.Client) *ioc.App {
	v := ioc2.InitDB()
	notificationDAO := dao.NewNotificationDAO(v)
	cmdable := ioc2.InitRedisCmd()
	quotaCache := redis.NewQuotaCache(cmdable)
	notificationRepository := repository.NewNotificationRepository(notificationDAO, quotaCache)
	service := notification.NewNotificationService(notificationRepository)
	channelTemplateDAO := dao.NewChannelTemplateDAO(v)
	channelTemplateRepository := repository.NewChannelTemplateRepository(channelTemplateDAO)
	string2 := ioc2.InitProviderEncryptKey()
	providerDAO := dao.NewProviderDAO(v, string2)
	providerRepository := repository.NewProviderRepository(providerDAO)
	manageService := manage.NewProviderService(providerRepository)
	auditService := audit.NewService()
	channelTemplateService := manage2.NewChannelTemplateService(channelTemplateRepository, manageService, auditService, clients)
	businessConfigDAO := dao.NewBusinessConfigDAO(v)
	redisClient := ioc2.InitRedisClient()
	cache := ioc2.InitGoCache()
	localCache := local.NewLocalCache(redisClient, cache)
	redisCache := redis.NewCache(redisClient)
	businessConfigRepository := repository.NewBusinessConfigRepository(businessConfigDAO, localCache, redisCache)
	businessConfigService := config.NewBusinessConfigService(businessConfigRepository)
	callbackLogDAO := dao.NewCallbackLogDAO(v)
	callbackLogRepository := repository.NewCallbackLogRepository(notificationRepository, callbackLogDAO)
	callbackService := callback.NewService(businessConfigService, callbackLogRepository)
	channel := newChannel(channelTemplateService, clients)
	taskPool := newTaskPool()
	notificationSender := sender.NewSender(notificationRepository, businessConfigService, callbackService, channel, taskPool)
	immediateSendStrategy := sendstrategy.NewImmediateStrategy(notificationRepository, notificationSender)
	defaultSendStrategy := sendstrategy.NewDefaultStrategy(notificationRepository, businessConfigService)
	sendStrategy := sendstrategy.NewDispatcher(immediateSendStrategy, defaultSendStrategy)
	sendService := notification.NewSendService(channelTemplateService, service, sendStrategy)
	txNotificationDAO := dao.NewTxNotificationDAO(v)
	txNotificationRepository := repository.NewTxNotificationRepository(txNotificationDAO)
	dlockClient := ioc2.InitDistributedLock(redisClient)
	txNotificationService := notification.NewTxNotificationService(txNotificationRepository, businessConfigService, notificationRepository, dlockClient, notificationSender)
	notificationServer := grpc.NewServer(service, sendService, txNotificationService, channelTemplateService)
	component := ioc2.InitEtcdClient()
	egrpcComponent := ioc2.InitGrpc(notificationServer, component)
	asyncRequestResultCallbackTask := callback.NewAsyncRequestResultCallbackTask(dlockClient, callbackService)
	notificationScheduler := scheduler.NewScheduler(service, notificationSender, dlockClient)
	sendingTimeoutTask := notification.NewSendingTimeoutTask(dlockClient, notificationRepository)
	txCheckTask := notification.NewTxCheckTask(txNotificationRepository, businessConfigService, dlockClient)
	v2 := ioc2.InitTasks(asyncRequestResultCallbackTask, notificationScheduler, sendingTimeoutTask, txCheckTask)
	quotaDAO := dao.NewQuotaDAO(v)
	quotaRepository := repository.NewQuotaRepository(quotaDAO)
	quotaService := quota.NewService(quotaRepository)
	monthlyResetCron := quota.NewQuotaMonthlyResetCron(businessConfigRepository, quotaService)
	v3 := ioc2.Crons(monthlyResetCron, businessConfigRepository)
	app := &ioc.App{
		GrpcServer:          egrpcComponent,
		Tasks:               v2,
		Crons:               v3,
		CallbackSvc:         callbackService,
		CallbackLogRepo:     callbackLogRepository,
		ConfigSvc:           businessConfigService,
		ConfigRepo:          businessConfigRepository,
		NotificationSvc:     service,
		SendNotificationSvc: sendService,
		NotificationRepo:    notificationRepository,
		ProviderSvc:         manageService,
		ProviderRepo:        providerRepository,
		QuotaSvc:            quotaService,
		QuotaRepo:           quotaRepository,
		TemplateSvc:         channelTemplateService,
		TemplateRepo:        channelTemplateRepository,
		TxNotificationSvc:   txNotificationService,
		TxNotificationRepo:  txNotificationRepository,
	}
	return app
}

// wire.go:

var (
	BaseSet              = wire.NewSet(ioc2.InitDB, ioc2.InitDistributedLock, ioc2.InitEtcdClient, ioc2.InitIDGenerator, ioc2.InitRedisClient, ioc2.InitGoCache, ioc2.InitRedisCmd, local.NewLocalCache, redis.NewCache)
	configSvcSet         = wire.NewSet(config.NewBusinessConfigService, repository.NewBusinessConfigRepository, dao.NewBusinessConfigDAO)
	notificationSvcSet   = wire.NewSet(redis.NewQuotaCache, notification.NewNotificationService, repository.NewNotificationRepository, dao.NewNotificationDAO, notification.NewSendingTimeoutTask)
	txNotificationSvcSet = wire.NewSet(notification.NewTxNotificationService, repository.NewTxNotificationRepository, dao.NewTxNotificationDAO, notification.NewTxCheckTask)
	senderSvcSet         = wire.NewSet(
		newChannel,
		newTaskPool, sender.NewSender,
	)
	sendNotificationSvcSet = wire.NewSet(notification.NewSendService, sendstrategy.NewDispatcher, sendstrategy.NewImmediateStrategy, sendstrategy.NewDefaultStrategy)
	callbackSvcSet         = wire.NewSet(callback.NewService, repository.NewCallbackLogRepository, dao.NewCallbackLogDAO, callback.NewAsyncRequestResultCallbackTask)
	providerSvcSet         = wire.NewSet(manage.NewProviderService, repository.NewProviderRepository, dao.NewProviderDAO, ioc2.InitProviderEncryptKey)
	templateSvcSet         = wire.NewSet(manage2.NewChannelTemplateService, repository.NewChannelTemplateRepository, dao.NewChannelTemplateDAO)
	schedulerSet           = wire.NewSet(scheduler.NewScheduler)
	quotaSvcSet            = wire.NewSet(quota.NewService, quota.NewQuotaMonthlyResetCron, repository.NewQuotaRepository, dao.NewQuotaDAO)
)

func newTaskPool() pool.TaskPool {
	type Config struct {
		InitGo           int           `yaml:"initGo"`
		CoreGo           int32         `yaml:"coreGo"`
		MaxGo            int32         `yaml:"maxGo"`
		MaxIdleTime      time.Duration `yaml:"maxIdleTime"`
		QueueSize        int           `yaml:"queueSize"`
		QueueBacklogRate float64       `yaml:"queueBacklogRate"`
	}
	var cfg Config
	if err := econf.UnmarshalKey("pool", &cfg); err != nil {
		panic(err)
	}
	p, err := pool.NewOnDemandBlockTaskPool(cfg.InitGo, cfg.QueueSize, pool.WithQueueBacklogRate(cfg.QueueBacklogRate), pool.WithMaxIdleTime(cfg.MaxIdleTime), pool.WithCoreGo(cfg.CoreGo), pool.WithMaxGo(cfg.MaxGo))
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	return p
}

func newChannel(
	templateSvc manage2.ChannelTemplateService,
	clients map[string]client.Client,
) channel.Channel {
	return channel.NewDispatcher(map[domain.Channel]channel.Channel{domain.ChannelSMS: channel.NewSMSChannel(newSMSSelectorBuilder(templateSvc, clients))})
}

func newSMSSelectorBuilder(
	templateSvc manage2.ChannelTemplateService,
	clients map[string]client.Client,
) *sequential.SelectorBuilder {
	providers := make([]provider.Provider, 0, len(clients))
	for k := range clients {
		providers = append(providers, sms.NewSMSProvider(
			k,
			templateSvc,
			clients[k],
		))
	}
	return sequential.NewSelectorBuilder(providers)
}
