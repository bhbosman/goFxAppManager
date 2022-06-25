package internal

import (
	"context"
	model2 "github.com/bhbosman/goFxAppManager/Serivce/model"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"reflect"
)

type fxApplicationInformation struct {
	ServiceId         model.ServiceIdentifier
	ServiceDependency model.ServiceIdentifier
	isDirty           bool
	Name              string
	Callback          messages.CreateAppCallbackFn
}

func newFxApplicationInformation(
	ServiceId model.ServiceIdentifier,
	ServiceDependency model.ServiceIdentifier,
	name string, callback messages.CreateAppCallbackFn) *fxApplicationInformation {
	return &fxApplicationInformation{
		ServiceId:         ServiceId,
		ServiceDependency: ServiceDependency,
		isDirty:           true,
		Name:              name,
		Callback:          callback,
	}
}

type Data struct {
	isDirty                 bool
	appContext              context.Context
	pubSub                  *pubsub.PubSub
	fxCreateAppsCallbackMap map[string]*fxApplicationInformation
	fxAppsMap               map[string]*fx.App
	messageRouter           *messageRouter.MessageRouter
	logger                  *zap.Logger
}

func NewData(
	applicationContext context.Context,
	FnApps []messages.CreateAppCallback,
	pubSub *pubsub.PubSub,
	logger *zap.Logger) (*Data, error) {
	result := &Data{
		appContext:              applicationContext,
		pubSub:                  pubSub,
		fxCreateAppsCallbackMap: make(map[string]*fxApplicationInformation),
		fxAppsMap:               make(map[string]*fx.App),
		messageRouter:           messageRouter.NewMessageRouter(),
		logger:                  logger,
	}
	result.messageRouter.Add(result.handleEmptyQueue)

	for _, app := range FnApps {
		result.fxCreateAppsCallbackMap[app.Name] = newFxApplicationInformation(
			app.ServiceId,
			app.ServiceDependency,
			app.Name,
			app.Callback)
		result.isDirty = true
		result.publish(
			&model2.FxServiceAdded{
				Name: app.Name,
			})
	}

	return result, nil
}

func (self *Data) StopAll(ctx context.Context) error {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "StopAll"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	var err error
	self.logger.Error("Starting all services",
		zap.String("Method", "StartAll"),
		zap.Error(self.appContext.Err()))
	for name := range self.fxCreateAppsCallbackMap {
		err = multierr.Append(err, self.Stop(ctx, name))
	}
	return err
}

func (self *Data) StartAll(startContext context.Context) error {
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

func (self *Data) Stop(stopContext context.Context, name ...string) error {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "Stop"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	var err error
	err = nil
	for _, iterName := range name {
		// check if not in started list
		var ok bool
		var app *fx.App
		if app, ok = self.fxAppsMap[iterName]; ok {
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

func (self *Data) Start(startContext context.Context, name ...string) error {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "Start"),
			zap.Error(self.appContext.Err()),
			zap.Strings("name", name))
		return self.appContext.Err()
	}
	var err error
	err = nil
	self.logger.Info("Starting service", zap.String("Method", "Start"), zap.Error(self.appContext.Err()), zap.Strings("name", name))
	for _, iterName := range name {
		// check if not in started list

		var ok bool
		var applicationInformation *fxApplicationInformation
		var app *fx.App
		var cancelFunc context.CancelFunc
		self.logger.Info("Check if already started", zap.String("ServiceName", iterName))
		if _, ok = self.fxAppsMap[iterName]; !ok {
			self.logger.Info("Not started", zap.String("ServiceName", iterName))
			if applicationInformation, ok = self.fxCreateAppsCallbackMap[iterName]; ok {
				self.logger.Info("Starting", zap.String("ServiceName", iterName))
				app, cancelFunc, err = applicationInformation.Callback()
				onError := func() {
					if cancelFunc != nil {
						cancelFunc()
					}
				}

				if err == nil {
					self.publish(&model2.FxServiceStarted{
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

func (self *Data) ShutDown() error {
	return self.appContext.Err()
}

func (self *Data) Send(message interface{}) error {
	if self.appContext.Err() != nil {
		self.logger.Error("App Context in Error",
			zap.String("Method", "Send"),
			zap.Error(self.appContext.Err()),
			zap.String("MessageType", reflect.TypeOf(message).String()),
		)
		return self.appContext.Err()
	}
	b, _ := self.messageRouter.Route(message)
	if !b {
		self.logger.Info("No message handler implemented for", zap.String("TypeName", reflect.TypeOf(message).String()))
	}
	return nil
}

func (self *Data) publish(message interface{}) {
	self.pubSub.Pub(message, "ActiveFxServicesStatus")
}

func (self *Data) handleEmptyQueue(message *messages.EmptyQueue) interface{} {
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
					&model2.FxServiceStatus{
						Name:              information.Name,
						Active:            active,
						ServiceId:         information.ServiceId,
						ServiceDependency: information.ServiceDependency,
					})
				information.isDirty = false
			}
		}
		self.isDirty = false
	}
	return nil
}
