package FxServicesSlide

import (
	"context"
	"github.com/bhbosman/goFxAppManager/Serivce"
	"github.com/bhbosman/goUi/ui"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
)

type Factory struct {
	fxManagerService   Serivce.IFxManagerService
	applicationContext context.Context
	pubSub             *pubsub.PubSub
	app                *tview.Application
}

func (self *Factory) OrderNumber() int {
	return 200
}

func NewFactory(
	applicationContext context.Context,
	pubSub *pubsub.PubSub,
	app *tview.Application,
	fxManagerService Serivce.IFxManagerService) (*Factory, error) {
	return &Factory{
		fxManagerService:   fxManagerService,
		applicationContext: applicationContext,
		pubSub:             pubSub,
		app:                app,
	}, nil
}

func (self *Factory) Title() string {
	return "FxServices"
}

func (self *Factory) Content() ui.SlideCallback {
	return func(
		nextSlide func()) (string, ui.IPrimitiveCloser) {
		return self.Title(),
			NewFxServiceSlide(
				self.applicationContext,
				self.pubSub, self.app,
				self.fxManagerService,
			)

	}
}
