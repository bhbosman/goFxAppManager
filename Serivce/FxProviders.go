package Serivce

import (
	"context"
	"github.com/bhbosman/goFxAppManager/Serivce/internal"
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
					OnStart: params.FxManagerService.OnStart,
					OnStop:  params.FxManagerService.OnStop,
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
						OnData             internal.OnDataCallback
						ApplicationContext context.Context `name:"Application"`
						PubSub             *pubsub.PubSub  `name:"Application"`
					},
				) (IFxManagerService, error) {
					return internal.NewFxManagerService(
						params.ApplicationContext,
						params.OnData,
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
				) internal.OnDataCallback {
					return func(applicationContext context.Context) (internal.IFxManagerData, error) {
						return internal.NewData(
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
