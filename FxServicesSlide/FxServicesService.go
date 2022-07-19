package FxServicesSlide

import (
	"context"
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/uiCommon"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Service struct {
	ctx                      context.Context
	cancelFunc               context.CancelFunc
	channel                  chan interface{}
	OnData                   func() (internal.IFxManagerData, error)
	state                    IFxService.State
	pubSub                   *pubsub.PubSub
	connectionListChange     func(connectionList []internal.IdAndName)
	connectionInstanceChange func(data internal.SendActionsForService)
	logger                   *zap.Logger
	goFunctionCounter        GoFunctionCounter.IService
	pubSubChannel            chan interface{}
}

func (self *Service) Send(message interface{}) error {
	send, err := ISendMessage.CallISendMessageSend(self.ctx, self.channel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *Service) StartService(name string) {
	_, _ = internal.CallIFxManagerSlideStartService(self.ctx, self.channel, false, name)
}

func (self *Service) StopService(name string) {
	_, _ = internal.CallIFxManagerSlideStopService(self.ctx, self.channel, false, name)
}

func (self *Service) StartAllService() {
	_, _ = internal.CallIFxManagerSlideStartAllService(self.ctx, self.channel, false)
}

func (self *Service) StopAllService() {
	_, _ = internal.CallIFxManagerSlideStopAllService(self.ctx, self.channel, false)
}

func (self *Service) SetConnectionListChange(cb func(connectionList []internal.IdAndName)) {
	self.connectionListChange = cb
}

func (self *Service) SetConnectionInstanceChange(cb func(data internal.SendActionsForService)) {
	self.connectionInstanceChange = cb
}

func (self *Service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *Service) start(_ context.Context) error {
	data, err := self.OnData()
	if err != nil {
		return err
	}
	data.SetConnectionListChange(self.connectionListChange)
	data.SetConnectionInstanceChange(self.connectionInstanceChange)

	return self.goFunctionCounter.GoRun("FxServiceSlide.Start",
		func() {
			self.goStart(data)
		},
	)
}
func (self *Service) goStart(data internal.IFxManagerData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.channel)

	self.pubSubChannel = self.pubSub.Sub("ActiveFxServicesStatus", uiCommon.UIState)

	var messageReceived interface{}
	var ok bool

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		data,
		[]ChannelHandler.ChannelHandler{
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(internal.IFxManagerSlide); ok {
						return internal.ChannelEventsForIFxManagerSlide(unk, message)
					}
					return false, nil
				},
			},
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return ISendMessage.ChannelEventsForISendMessage(unk, message)
					}
					return false, nil
				},
			},
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if sm, ok := next.(ISendMessage.ISendMessage); ok {
						_ = sm.Send(message)
					}
					return true, nil
				},
			},
		},
		func() int {
			n := len(self.pubSubChannel) + len(self.channel)
			return n
		},
		//func(i interface{}) {
		//	select {
		//	case self.channel <- i:
		//		break
		//	default:
		//		break
		//	}
		//},
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := data.ShutDown()
			if err != nil {
				self.logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case messageReceived, ok = <-self.channel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived)
			if err != nil || b {
				return
			}
			break
		case messageReceived, ok = <-self.pubSubChannel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived)
			if err != nil || b {
				return
			}
			break
		}
	}
}

func (self *Service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	close(self.channel)
	self.state = IFxService.Stopped
	return err
}

func (self *Service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("FxAppManager.PubSub.Unsubscribe", self.pubSub, self.goFunctionCounter, self.pubSubChannel)
}

func (self *Service) State() IFxService.State {
	return self.state
}

func (self *Service) ServiceName() string {
	return "FxServicesSlide"
}

func NewService(
	applicationContext context.Context,
	OnData func() (internal.IFxManagerData, error),
	pubSub *pubsub.PubSub,
	logger *zap.Logger,
	goFunctionCounter GoFunctionCounter.IService,
) (internal.IFxManagerService, error) {
	ctx, cancelFunc := context.WithCancel(applicationContext)
	channel := make(chan interface{}, 32)
	return &Service{
		ctx:               ctx,
		cancelFunc:        cancelFunc,
		channel:           channel,
		OnData:            OnData,
		pubSub:            pubSub,
		logger:            logger,
		goFunctionCounter: goFunctionCounter,
	}, nil
}
