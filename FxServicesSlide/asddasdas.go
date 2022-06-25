package FxServicesSlide

import (
	"context"
	"github.com/bhbosman/goFxAppManager/Serivce"
	"github.com/bhbosman/goUi/ui"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

func Dddddd() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "RegisteredMainWindowSlides",
			Target: func(
				params struct {
					fx.In
					App                *tview.Application
					ApplicationContext context.Context `name:"Application"`
					PubSub             *pubsub.PubSub  `name:"Application"`
					FxManagerService   Serivce.IFxManagerService
				},
			) (ui.ISlideFactory, error) {
				return NewFactory(
					params.ApplicationContext,
					params.PubSub,
					params.App,
					params.FxManagerService)
			}})

}
