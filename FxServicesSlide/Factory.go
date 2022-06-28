package FxServicesSlide

import (
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
	"github.com/bhbosman/goUi/ui"
	"github.com/rivo/tview"
)

type Factory struct {
	app     *tview.Application
	service internal.IFxManagerService
}

func (self *Factory) OrderNumber() int {
	return 200
}

func NewFactory(
	app *tview.Application,
	service internal.IFxManagerService,
) (*Factory, error) {
	return &Factory{
		app:     app,
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
