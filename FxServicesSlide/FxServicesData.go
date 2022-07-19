package FxServicesSlide

import (
	"context"
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
	"github.com/bhbosman/goFxAppManager/service"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/uiCommon"
	"sort"
)

type fxServiceManagerSlideData struct {
	serviceListIsDirty         bool
	ss                         map[string]*FxServicesManagerData
	isDirty                    map[string]bool
	messageRouter              *messageRouter.MessageRouter
	onConnectionListChange     func(connectionList []internal.IdAndName)
	onConnectionInstanceChange func(data internal.SendActionsForService)
	fxManagerService           service.IFxManagerService
	UiActive                   bool
}

func (self *fxServiceManagerSlideData) ShutDown() error {
	return nil
}

func (self *fxServiceManagerSlideData) StartAllService() {
	_ = self.fxManagerService.StartAll(context.Background())
}

func (self *fxServiceManagerSlideData) StopAllService() {
	_ = self.fxManagerService.StopAll(context.Background())
}

func (self *fxServiceManagerSlideData) StartService(name string) {
	_ = self.fxManagerService.Start(context.Background(), name)
}

func (self *fxServiceManagerSlideData) StopService(name string) {
	_ = self.fxManagerService.Stop(context.Background(), name)
}

func (self *fxServiceManagerSlideData) Send(data interface{}) error {
	_, err := self.messageRouter.Route(data)
	return err
}
func (self *fxServiceManagerSlideData) handleUiStarted(message *uiCommon.UiStarted) {
	self.UiActive = message.Active
}

func (self *fxServiceManagerSlideData) handlePublishInstanceDataFor(message *publishInstanceDataFor) error {
	self.isDirty[message.Name] = true
	return nil
}

func (self *fxServiceManagerSlideData) handleEmptyQueue(_ *messages.EmptyQueue) {
	if self.UiActive {
		if self.serviceListIsDirty {
			self.DoServiceListChange()
			self.serviceListIsDirty = false
		}
		for key := range self.isDirty {
			if v, ok := self.ss[key]; ok {
				self.DoServiceInstanceChange(v)
			}
		}
		self.isDirty = make(map[string]bool)

	}
}

func (self *fxServiceManagerSlideData) DoServiceListChange() {
	if self.onConnectionListChange != nil {
		ss := make([]string, 0, len(self.ss))

		for key := range self.ss {
			ss = append(ss, key)
		}
		sort.Strings(ss)
		cbData := make([]internal.IdAndName, 0, len(self.ss))

		for _, s := range ss {
			if info, ok := self.ss[s]; ok {
				idAndName := internal.IdAndName{
					ServiceId:         info.ServiceId,
					ServiceDependency: info.ServiceDependency,
					Name:              info.Name,
					Active:            info.Active,
				}
				cbData = append(cbData, idAndName)
			}
		}
		self.onConnectionListChange(cbData)
	}
}

const StopServiceString = "Stop Service"
const StartServiceString = "Start Service"
const StopAllServiceString = "Stop All Service"
const StartAllServiceString = "Start All Service"

func (self *fxServiceManagerSlideData) DoServiceInstanceChange(data *FxServicesManagerData) {
	if self.onConnectionInstanceChange != nil {
		var actions []string
		if data.Active {
			actions = append(actions, StopServiceString)
		} else {
			actions = append(actions, StartServiceString)
		}
		actions = append(actions, []string{"-", StartAllServiceString, StopAllServiceString}...)
		SendActionsForService := internal.SendActionsForService{
			Name:    data.Name,
			Actions: actions,
		}
		dd := SendActionsForService
		self.onConnectionInstanceChange(dd)
	}
}

func (self *fxServiceManagerSlideData) SetConnectionInstanceChange(cb func(data internal.SendActionsForService)) {
	self.onConnectionInstanceChange = cb
}

func (self *fxServiceManagerSlideData) SetConnectionListChange(cb func(connectionList []internal.IdAndName)) {
	self.onConnectionListChange = cb
}

func (self *fxServiceManagerSlideData) handleFxServiceStatus(message *service.FxServiceStatus) error {
	self.ss[message.Name] = &FxServicesManagerData{
		Name:              message.Name,
		Active:            message.Active,
		ServiceId:         message.ServiceId,
		ServiceDependency: message.ServiceDependency,
	}
	self.isDirty[message.Name] = true
	self.serviceListIsDirty = true
	return nil
}

func (self *fxServiceManagerSlideData) handleFxServiceStarted(message *service.FxServiceStarted) error {
	if instance, ok := self.ss[message.Name]; ok {
		instance.Active = true
		self.isDirty[message.Name] = true

	}
	return nil
}

func (self *fxServiceManagerSlideData) handleFxServiceStopped(message *service.FxServiceStopped) error {
	if instance, ok := self.ss[message.Name]; ok {
		instance.Active = false
		self.isDirty[message.Name] = true

	}
	return nil
}

func (self *fxServiceManagerSlideData) handleFxServiceAdded(message *service.FxServiceAdded) error {
	self.ss[message.Name] = &FxServicesManagerData{
		Name:   message.Name,
		Active: false,
	}
	self.isDirty[message.Name] = true
	self.serviceListIsDirty = true

	return nil
}

func NewData(
	fxManagerService service.IFxManagerService,
) (*fxServiceManagerSlideData, error) {
	result := &fxServiceManagerSlideData{
		ss:               make(map[string]*FxServicesManagerData),
		isDirty:          make(map[string]bool),
		messageRouter:    messageRouter.NewMessageRouter(),
		fxManagerService: fxManagerService,
	}
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleFxServiceStarted)
	_ = result.messageRouter.Add(result.handleFxServiceStopped)
	_ = result.messageRouter.Add(result.handleFxServiceAdded)
	_ = result.messageRouter.Add(result.handleFxServiceStatus)
	_ = result.messageRouter.Add(result.handlePublishInstanceDataFor)
	_ = result.messageRouter.Add(result.handleUiStarted)

	return result, nil
}
