package internal

import (
	"context"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"go.uber.org/zap"
)

type OnDataCallback func(applicationContext context.Context) (IFxManagerData, error)

type Service struct {
	context           context.Context
	cancelFunc        context.CancelFunc
	channel           chan interface{}
	onData            OnDataCallback
	state             IFxService.State
	logger            *zap.Logger
	goFunctionCounter GoFunctionCounter.IService
}

func (self *Service) ServiceName() string {
	return "FxAppManagerService"
}

func (self *Service) State() IFxService.State {
	return self.state
}

func NewFxManagerService(
	applicationContext context.Context,
	onData OnDataCallback,
	logger *zap.Logger,
	goFunctionCounter GoFunctionCounter.IService,
) (*Service, error) {
	ctx, cancelFunc := context.WithCancel(applicationContext)
	return &Service{
		context:           ctx,
		cancelFunc:        cancelFunc,
		channel:           make(chan interface{}, 32),
		onData:            onData,
		state:             IFxService.NotInitialized,
		logger:            logger,
		goFunctionCounter: goFunctionCounter,
	}, nil
}

func (self *Service) StopAll(ctx context.Context) error {
	result, err := CallIFxManagerStopAll(self.context, self.channel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) StartAll(ctx context.Context) error {
	result, err := CallIFxManagerStartAll(self.context, self.channel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) StopStartAll(ctx context.Context) error {
	result, err := CallIFxManagerStartAll(self.context, self.channel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) Stop(ctx context.Context, name ...string) error {
	result, err := CallIFxManagerStop(self.context, self.channel, true, ctx, name...)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) Start(ctx context.Context, name ...string) error {
	result, err := CallIFxManagerStart(self.context, self.channel, true, ctx, name...)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) OnStart(ctx context.Context) error {
	err := self.start()
	if err != nil {
		return err
	}

	err = self.StartAll(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *Service) OnStop(_ context.Context) error {
	err := self.shutdown()
	close(self.channel)
	self.state = IFxService.Stopped
	return err
}

func (self *Service) start() error {
	data, err := self.onData(self.context)
	if err != nil {
		return err
	}

	// this function is part of the GoFunctionCounter count
	go func() {
		functionName := self.goFunctionCounter.CreateFunctionName("FxAppManager.start")
		defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
			_ = GoFunctionCounter.Remove(name)
		}(self.goFunctionCounter, functionName)
		_ = self.goFunctionCounter.Add(functionName)

		//
		self.goStart(data)
	}()
	return nil
}

func (self *Service) shutdown() error {
	self.cancelFunc()
	return nil
}

func (self *Service) goStart(data IFxManagerData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.channel)

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.context,
		data, []ChannelHandler.ChannelHandler{
			{
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IFxManager); ok {
						return ChannelEventsForIFxManager(unk, message)
					}
					return false, nil
				},
			},
			{
				BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return ISendMessage.ChannelEventsForISendMessage(unk, message)
					}
					return false, nil
				},
			},
		},
		func() int {
			return len(self.channel)
		})

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
		case event, ok := <-self.channel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break
			}
		}
	}
}
