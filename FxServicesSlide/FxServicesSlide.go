package FxServicesSlide

import (
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
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
}

func (self *FxServicesManagerSlide) UpdateContent() error {
	return nil
}

func (self *FxServicesManagerSlide) Close() error {
	//self.cancelFunc()
	//close(self.channel)
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

//func (self *FxServicesManagerSlide) goRun() {
//	defer func(cmdChannel <-chan interface{}) {
//		//flush
//		for range cmdChannel {
//		}
//	}(self.channel)
//
//	pubSubChannel := self.pubSub.Sub("ActiveFxServicesStatus")
//	defer func(pubSubChannel chan interface{}) {
//		// unsubscribe on different go routine to avoid deadlock
//		go func(pubSubChannel chan interface{}) {
//			self.pubSub.Unsub(pubSubChannel)
//			for range pubSubChannel {
//			}
//		}(pubSubChannel)
//	}(pubSubChannel)
//
//	var messageReceived interface{}
//	var ok bool
//
//	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
//		self.ctx,
//		self.data,
//		[]ChannelHandler.ChannelHandler{
//			{
//				BreakOnSuccess: false,
//				Cb: func(next interface{}, message interface{}) (bool, error) {
//					if unk, ok := next.(internal.IFxManagerSlide); ok {
//						return internal.ChannelEventsForIFxManagerSlide(unk, message)
//					}
//					return false, nil
//				},
//			},
//			{
//				PubSubHandler:  false,
//				BreakOnSuccess: false,
//				Cb: func(next interface{}, message interface{}) (bool, error) {
//					if unk, ok := next.(ISendMessage.ISendMessage); ok {
//						return ISendMessage.ChannelEventsForISendMessage(unk, message)
//					}
//					return false, nil
//				},
//			},
//			{
//				PubSubHandler:  true,
//				BreakOnSuccess: false,
//				Cb: func(next interface{}, message interface{}) (bool, error) {
//					if sm, ok := next.(ISendMessage.ISendMessage); ok {
//						_ = sm.Send(message)
//					}
//					return true, nil
//				},
//			},
//		},
//		func() int {
//			n := len(pubSubChannel) + len(self.channel)
//			return n
//		})
//loop:
//	for {
//		select {
//		case <-self.ctx.Done():
//			break loop
//		case messageReceived, ok = <-self.channel:
//			if !ok {
//				return
//			}
//			b, err := channelHandlerCallback(messageReceived, false)
//			if err != nil || b {
//				return
//			}
//			break
//		case messageReceived, ok = <-pubSubChannel:
//			if !ok {
//				return
//			}
//			b, err := channelHandlerCallback(messageReceived, true)
//			if err != nil || b {
//				return
//			}
//			break
//		}
//	}
//}

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
						self.service.StartService(data.Name)
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
					&PublishInstanceDataFor{
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
) *FxServicesManagerSlide {

	result := &FxServicesManagerSlide{
		service: service,
		app:     app,
	}
	result.init()
	service.SetConnectionListChange(result.SetFxServiceListChange)
	service.SetConnectionInstanceChange(result.SetFxServiceInstanceChange)
	return result
}
