package service

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type OnDataCallback func(applicationContext context.Context) (IFxManagerData, error)

type service struct {
	context           context.Context
	cancelFunc        context.CancelFunc
	cmdChannel        chan interface{}
	onData            OnDataCallback
	state             IFxService.State
	logger            *zap.Logger
	goFunctionCounter GoFunctionCounter.IService
}

func (self *service) GetState() (started []string, err error) {
	state, err := CallIFxManagerGetState(self.context, self.cmdChannel, true)
	if err != nil {
		return nil, err
	}
	return state.Args0, state.Args1
}

func (self *service) Add(name string, callback messages.CreateAppCallbackFn) error {
	add, err := CallIFxManagerAdd(self.context, self.cmdChannel, true, name, callback)
	if err != nil {
		return err
	}
	return add.Args0
}

func (self *service) Send(message interface{}) error {
	send, err := CallIFxManagerSend(self.context, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) ServiceName() string {
	return "FxAppManagerService"
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self *service) StopAll(ctx context.Context) error {
	result, err := CallIFxManagerStopAll(self.context, self.cmdChannel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) StartAll(ctx context.Context) error {
	result, err := CallIFxManagerStartAll(self.context, self.cmdChannel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) StopStartAll(ctx context.Context) error {
	result, err := CallIFxManagerStartAll(self.context, self.cmdChannel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) Stop(ctx context.Context, name ...string) error {
	result, err := CallIFxManagerStop(self.context, self.cmdChannel, true, ctx, name...)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) Start(ctx context.Context, name ...string) error {
	result, err := CallIFxManagerStart(self.context, self.cmdChannel, true, ctx, name...)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) OnStart(ctx context.Context) error {
	err := self.start()
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) OnStop(ctx context.Context) error {
	err := self.StopAll(ctx)
	//err = multierr.Append(err, self.closeAll())
	err = multierr.Append(err, self.shutdown())
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *service) start() error {
	data, err := self.onData(self.context)
	if err != nil {
		return err
	}

	return self.goFunctionCounter.GoRun("FxAppManager.start",
		func() {
			self.goStart(data)
		},
	)
}

func (self *service) shutdown() error {
	self.cancelFunc()
	return nil
}

func (self *service) goStart(data IFxManagerData) {
	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.context,
		data, []ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IFxManager); ok {
						return ChannelEventsForIFxManager(unk, message)
					}
					return false, nil
				},
			},
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return ISendMessage.ChannelEventsForISendMessage(unk, message)
					}
					return false, nil
				},
			},
		},
		func() int {
			return len(self.cmdChannel)
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
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
		case <-self.context.Done():
			err := data.ShutDown()
			if err != nil {
				self.logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case event, ok := <-self.cmdChannel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break
			}
		}
	}
	// flush
	for range self.cmdChannel {
	}
}

func newService(
	applicationContext context.Context,
	onData OnDataCallback,
	logger *zap.Logger,
	goFunctionCounter GoFunctionCounter.IService,
) (*service, error) {
	ctx, cancelFunc := context.WithCancel(applicationContext)
	return &service{
		context:           ctx,
		cancelFunc:        cancelFunc,
		cmdChannel:        make(chan interface{}, 32),
		onData:            onData,
		state:             IFxService.NotInitialized,
		logger:            logger,
		goFunctionCounter: goFunctionCounter,
	}, nil
}
