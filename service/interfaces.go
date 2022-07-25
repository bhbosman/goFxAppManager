package service

import (
	"context"
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type IFxManager interface {
	ISendMessage.ISendMessage
	Add(name string,
		callback messages.CreateAppCallbackFn,
	) error
	StopAll(ctx context.Context) error
	StartAll(ctx context.Context) error
	Stop(ctx context.Context, name ...string) error
	Start(ctx context.Context, name ...string) error
}

type IFxManagerService interface {
	IFxManager
	IFxService.IFxServices
}

type IFxManagerData interface {
	IFxManager
	IDataShutDown.IDataShutDown
}
