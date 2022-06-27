package FxServicesSlide

import (
	"context"
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/cskr/pubsub"
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
	go self.goStart(data)
	return nil
}
func (self *Service) goStart(data internal.IFxManagerData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.channel)

	pubSubChannel := self.pubSub.Sub("ActiveFxServicesStatus")
	defer func(pubSubChannel chan interface{}) {
		// unsubscribe on different go routine to avoid deadlock
		go func(pubSubChannel chan interface{}) {
			self.pubSub.Unsub(pubSubChannel)
			//flush
			for range pubSubChannel {
			}
		}(pubSubChannel)
	}(pubSubChannel)

	var messageReceived interface{}
	var ok bool

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		data,
		[]ChannelHandler.ChannelHandler{
			{
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(internal.IFxManagerSlide); ok {
						return internal.ChannelEventsForIFxManagerSlide(unk, message)
					}
					return false, nil
				},
			},
			{
				PubSubHandler:  false,
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return ISendMessage.ChannelEventsForISendMessage(unk, message)
					}
					return false, nil
				},
			},
			{
				PubSubHandler:  false,
				BreakOnSuccess: true,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IDataShutDown.IDataShutDown); ok {
						return IDataShutDown.ChannelEventsForIDataShutDown(unk, message)
					}
					return false, nil
				},
			},
			{
				PubSubHandler:  true,
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if sm, ok := next.(ISendMessage.ISendMessage); ok {
						_ = sm.Send(message)
					}
					return true, nil
				},
			},
		},
		func() int {
			n := len(pubSubChannel) + len(self.channel)
			return n
		})
loop:
	for {
		select {
		case <-self.ctx.Done():
			break loop
		case messageReceived, ok = <-self.channel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived, false)
			if err != nil || b {
				return
			}
			break
		case messageReceived, ok = <-pubSubChannel:
			if !ok {
				return
			}
			b, err := channelHandlerCallback(messageReceived, true)
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
	data, err := IDataShutDown.CallIDataShutDownShutDown(self.ctx, self.channel, true)
	if err != nil {
		return err
	}
	self.cancelFunc()
	return data.Args0
}

func (self *Service) State() IFxService.State {
	//TODO implement me
	panic("implement me")
}

func (self *Service) ServiceName() string {
	//TODO implement me
	panic("implement me")
}

func NewService(
	applicationContext context.Context,
	OnData func() (internal.IFxManagerData, error),
	pubSub *pubsub.PubSub,
) (internal.IFxManagerService, error) {
	ctx, cancelFunc := context.WithCancel(applicationContext)
	channel := make(chan interface{}, 32)
	return &Service{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		channel:    channel,
		OnData:     OnData,
		pubSub:     pubSub,
	}, nil
}
