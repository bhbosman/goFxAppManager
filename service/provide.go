package service

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func InvokeFxManager() fx.Option {
	return fx.Options(
		fx.Invoke(
			func(
				params struct {
					fx.In
					Lifecycle        fx.Lifecycle
					FxManagerService IFxManagerService
				},
			) error {
				hook := fx.Hook{
					OnStart: func(ctx context.Context) error {
						return params.FxManagerService.OnStart(ctx)
					},
					OnStop: params.FxManagerService.OnStop,
				}
				params.Lifecycle.Append(hook)
				return nil
			},
		),
	)
}

func ProvideFxManager() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						OnData             OnDataCallback
						ApplicationContext context.Context `name:"Application"`
						PubSub             *pubsub.PubSub  `name:"Application"`
						Logger             *zap.Logger
						GoFunctionCounter  GoFunctionCounter.IService
					},
				) (IFxManagerService, error) {
					return NewService(
						params.ApplicationContext,
						params.OnData,
						params.Logger,
						params.GoFunctionCounter,
					)
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						PubSub *pubsub.PubSub               `name:"Application"`
						FnApps []messages.CreateAppCallback `group:"Apps"`
						Logger *zap.Logger
					},
				) OnDataCallback {
					return func(applicationContext context.Context) (IFxManagerData, error) {
						return NewData(
							applicationContext,
							params.FnApps,
							params.PubSub,
							params.Logger.Named("FxServiceData"),
						)
					}
				},
			},
		),
	)
}
