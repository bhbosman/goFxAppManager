package FxServicesSlide

import (
	"context"
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
	"github.com/bhbosman/goFxAppManager/service"
	"github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideServiceSlide() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						FxManagerService service.IFxManagerService
					},
				) (func() (internal.IFxManagerData, error), error) {
					return func() (internal.IFxManagerData, error) {
						return NewData(params.FxManagerService)
					}, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						ApplicationContext context.Context `name:"Application"`
						PubSub             *pubsub.PubSub  `name:"Application"`
						Lifecycle          fx.Lifecycle
						Logger             *zap.Logger
						GoFunctionCounter  GoFunctionCounter.IService
						OnData             func() (internal.IFxManagerData, error)
					},
				) (internal.IFxManagerService, error) {
					service, err := NewService(
						params.ApplicationContext,
						params.OnData,
						params.PubSub,
						params.Logger,
						params.GoFunctionCounter,
					)
					if err != nil {
						return nil, err
					}
					params.Lifecycle.Append(
						fx.Hook{
							OnStart: func(ctx context.Context) error {
								return service.OnStart(ctx)
							},
							OnStop: func(ctx context.Context) error {
								return service.OnStop(ctx)
							},
						})
					return service, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Group: "RegisteredMainWindowSlides",
				Target: func(
					params struct {
						fx.In
						App     *tview.Application
						Service internal.IFxManagerService
					},
				) (ui.ISlideFactory, error) {
					return NewFactory(
						params.App,
						params.Service,
					)
				},
			},
		),
	)
}
