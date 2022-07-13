package FxServicesSlide

import (
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
	"github.com/bhbosman/goUi/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FxServicesManagerSlide struct {
	service    internal.IFxManagerService
	table      *tview.Table
	actionList *tview.List
	next       tview.Primitive
	app        *tview.Application
	plate      *PlateContent
	toggle     bool
}

func (self *FxServicesManagerSlide) Toggle(b bool) {
	self.toggle = b
}

func (self *FxServicesManagerSlide) UpdateContent() error {
	return nil
}

func (self *FxServicesManagerSlide) Close() error {
	return nil
}

func (self *FxServicesManagerSlide) Draw(screen tcell.Screen) {
	self.next.Draw(screen)
}

func (self *FxServicesManagerSlide) GetRect() (int, int, int, int) {
	return self.next.GetRect()
}

func (self *FxServicesManagerSlide) SetRect(x, y, width, height int) {
	self.next.SetRect(x, y, width, height)
}

func (self *FxServicesManagerSlide) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.next.InputHandler()
}

func (self *FxServicesManagerSlide) Focus(delegate func(p tview.Primitive)) {
	self.next.Focus(delegate)
}

func (self *FxServicesManagerSlide) HasFocus() bool {
	return self.next.HasFocus()
}

func (self *FxServicesManagerSlide) Blur() {
	self.next.Blur()
}

func (self *FxServicesManagerSlide) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.next.MouseHandler()
}

func (self *FxServicesManagerSlide) SetFxServiceListChange(list []internal.IdAndName) {
	self.app.QueueUpdateDraw(func() {
		self.plate = newFxAppManagerPlateContent(list)
		self.table.SetContent(self.plate)
	})
}

func (self *FxServicesManagerSlide) SetFxServiceInstanceChange(data internal.SendActionsForService) {
	self.app.QueueUpdateDraw(func() {
		self.actionList.Clear()
		self.actionList.AddItem("..", "", 0, func() {
			self.app.SetFocus(self.table)
		})
		for _, action := range data.Actions {
			if action == StopServiceString {
				self.actionList.AddItem(action, "", 0,
					func() {
						self.service.StopService(data.Name)
						self.app.SetFocus(self.table)
					},
				)
				continue
			}
			if action == StartServiceString {
				self.actionList.AddItem(action, "", 0,
					func() {
						self.service.StartService(data.Name)
						self.app.SetFocus(self.table)
					},
				)
				continue
			}
			if action == StartAllServiceString {
				self.actionList.AddItem(action, "", 0,
					func() {
						self.service.StartAllService()
						self.app.SetFocus(self.table)
					},
				)
				continue
			}
			if action == StopAllServiceString {
				self.actionList.AddItem(action, "", 0,
					func() {
						self.service.StopAllService()
						self.app.SetFocus(self.table)
					},
				)
				continue
			}
			self.actionList.AddItem(action, "", 0, nil)

		}
	})
}

func (self *FxServicesManagerSlide) init() {

	self.actionList = tview.NewList().ShowSecondaryText(false)
	self.actionList.SetBorder(true).SetTitle("Actions")
	self.table = tview.NewTable()
	self.table.
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetSelectedFunc(func(row, column int) {
			self.app.SetFocus(self.actionList)
		}).
		SetSelectionChangedFunc(
			func(row, column int) {
				_ = self.service.Send(
					&publishInstanceDataFor{
						Name: self.plate.Grid[row-1].Name,
					},
				)
			},
		).
		SetBorder(true).
		SetTitle("Service Manager")
	self.next = tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(self.table, 0, 3, true),
					0, 5, true).
				AddItem(self.actionList, 0, 1, false),
			0,
			1,
			true)
}

func NewFxServiceSlide(
	service internal.IFxManagerService,
	app *tview.Application,
) ui.IPrimitiveCloser {

	result := &FxServicesManagerSlide{
		service: service,
		app:     app,
	}
	result.init()
	service.SetConnectionListChange(result.SetFxServiceListChange)
	service.SetConnectionInstanceChange(result.SetFxServiceInstanceChange)
	return result
}
