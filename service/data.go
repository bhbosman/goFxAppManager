package service

import (
	"context"
	"github.com/bhbosman/goConn"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"reflect"
)

type FxApplicationInformation struct {
	Name     string
	Callback messages.CreateAppCallbackFn
	isDirty  bool
}

func NewFxApplicationInformation(
	name string,
	callback messages.CreateAppCallbackFn,
) *FxApplicationInformation {
	return &FxApplicationInformation{
		isDirty:  true,
		Name:     name,
		Callback: callback,
	}
}

type data struct {
	isDirty                 bool
	appContext              context.Context
	pubSub                  *pubsub.PubSub
	fxCreateAppsCallbackMap map[string]*FxApplicationInformation
	fxAppsMap               map[string]messages.IApp
	messageRouter           messageRouter.IMessageRouter
	logger                  *zap.Logger
}

func (self *data) Publish() error {
	self.isDirty = true
	return nil
}

func (self *data) GetState() (started []string, err error) {
	result := make([]string, 0, len(self.fxAppsMap))
	for key := range self.fxAppsMap {
		result = append(result, key)
	}
	return result, nil
}

func (self *data) SaveToConfiguration() error {
	return nil
}

func (self *data) StartFromConfiguration() error {
	return nil
}

func (self *data) Add(
	name string,
	callback messages.CreateAppCallbackFn,
) error {
	self.fxCreateAppsCallbackMap[name] = NewFxApplicationInformation(
		name,
		callback)
	self.isDirty = true
	self.publish(
		&FxServiceAdded{
			Name: name,
		},
	)
	return nil
}

func (self *data) StopAll(ctx context.Context) error {
	var err error
	for name := range self.fxCreateAppsCallbackMap {
		err = multierr.Append(err, self.Stop(ctx, name))
	}
	return err
}

func (self *data) StartAll(startContext context.Context) error {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "StartAll"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	var err error
	self.logger.Error("Starting all services",
		zap.String("Method", "StartAll"),
		zap.Error(self.appContext.Err()))
	for name := range self.fxCreateAppsCallbackMap {
		err = multierr.Append(err, self.Start(startContext, name))
	}
	return err
}

func (self *data) Stop(stopContext context.Context, name ...string) error {
	var err error
	for _, iterName := range name {
		var ok bool
		var app messages.IApp
		if app, ok = self.fxAppsMap[iterName]; ok {
			if app == nil {
				continue
			}
			err = app.Stop(stopContext)
			delete(self.fxAppsMap, iterName)
			if instance, ok := self.fxCreateAppsCallbackMap[iterName]; ok {
				instance.isDirty = true
			}
			self.isDirty = true
		}
	}
	return err
}

func (self *data) Start(startContext context.Context, name ...string) error {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "Start"),
			zap.Error(self.appContext.Err()),
			zap.Strings("name", name))
		return self.appContext.Err()
	}
	var err error
	err = nil
	self.logger.Info(
		"Starting service",
		zap.String("Method", "Start"),
		zap.Error(self.appContext.Err()),
		zap.Strings("name", name))
	for _, iterName := range name {
		// check if not in started list

		var ok bool
		var applicationInformation *FxApplicationInformation
		var app messages.IApp
		var cancelFunc goConn.ICancellationContext
		self.logger.Info("Check if already started", zap.String("ServiceName", iterName))
		if _, ok = self.fxAppsMap[iterName]; !ok {
			self.logger.Info("Not started", zap.String("ServiceName", iterName))
			if applicationInformation, ok = self.fxCreateAppsCallbackMap[iterName]; ok {
				self.logger.Info("Starting", zap.String("ServiceName", iterName))
				app, cancelFunc, err = applicationInformation.Callback()
				onError := func() {
					if cancelFunc != nil {
						cancelFunc.CancelFunc()
					}
				}

				if err == nil {
					self.publish(&FxServiceStarted{
						Name: iterName,
					})
					applicationInformation.isDirty = true
					self.isDirty = true
					// start service
					err = app.Start(startContext)
					if err == nil {
						// no error add to started list
						self.fxAppsMap[applicationInformation.Name] = app
					} else {
						self.logger.Error("Starting error", zap.String("ServiceName", iterName), zap.Error(err))
						err = multierr.Append(err, err)
						onError()
					}
				} else {
					self.logger.Error("Creation error", zap.String("ServiceName", iterName), zap.Error(err))
					err = multierr.Append(err, err)
					onError()
				}
			}
		}
	}
	return err
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) Send(message interface{}) error {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "Send"),
			zap.Error(self.appContext.Err()),
			zap.String("MessageType", reflect.TypeOf(message).String()),
		)
		return self.appContext.Err()
	}
	self.messageRouter.Route(message)
	return nil
}

func (self *data) publish(message interface{}) {
	self.pubSub.Pub(message, "ActiveFxServicesStatus")
}

func (self *data) handleEmptyQueue(message *messages.EmptyQueue) interface{} {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "Send"),
			zap.Error(self.appContext.Err()),
			zap.String("MessageType", reflect.TypeOf(message).String()),
		)
		return self.appContext.Err()
	}
	if self.isDirty {
		for _, information := range self.fxCreateAppsCallbackMap {
			if information.isDirty {
				_, active := self.fxAppsMap[information.Name]
				self.publish(
					&FxServiceStatus{
						Name:   information.Name,
						Active: active,
					},
				)
				information.isDirty = false
			}
		}
		self.isDirty = false
	}
	return nil
}

func newData(
	applicationContext context.Context,
	FnApps []messages.CreateAppCallback,
	pubSub *pubsub.PubSub,
	logger *zap.Logger) (*data, error) {
	result := &data{
		appContext:              applicationContext,
		pubSub:                  pubSub,
		fxCreateAppsCallbackMap: make(map[string]*FxApplicationInformation),
		fxAppsMap:               make(map[string]messages.IApp),
		messageRouter:           messageRouter.NewMessageRouter(),
		logger:                  logger,
	}
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	var err error
	for _, app := range FnApps {
		err = multierr.Append(
			err,
			result.Add(
				app.Name,
				app.Callback,
			),
		)
	}

	return result, nil
}
