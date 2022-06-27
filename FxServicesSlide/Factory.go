package FxServicesSlide

import (
	"context"
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
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
	//onData             func() (internal.IFxManagerData, error)
	service internal.IFxManagerService
}

func (self *Factory) OrderNumber() int {
	return 200
}

func NewFactory(
	applicationContext context.Context,
	pubSub *pubsub.PubSub,
	app *tview.Application,
	fxManagerService Serivce.IFxManagerService,
	//onData func() (internal.IFxManagerData, error),
	service internal.IFxManagerService,
) (*Factory, error) {
	return &Factory{
		fxManagerService:   fxManagerService,
		applicationContext: applicationContext,
		pubSub:             pubSub,
		app:                app,
		//onData:             onData,
		service: service,
	}, nil
}

func (self *Factory) Title() string {
	return "FxServices"
}

func (self *Factory) Content(nextSlide func()) (string, ui.IPrimitiveCloser, error) {
	slide := NewFxServiceSlide(
		self.service,
		self.app,
	)
	return self.Title(), slide, nil
}
