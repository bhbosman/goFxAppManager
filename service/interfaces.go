package service

import (
	"context"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/services/IDataShutDown"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

//gfdlgdfjlj

type IFxManager interface {
	ISendMessage.ISendMessage
	Add(name string, callback gocommon.CreateAppCallbackFn) error
	StopAll(ctx context.Context) error
	StartAll(ctx context.Context) error
	Stop(ctx context.Context, name ...string) error
	Start(ctx context.Context, name ...string) error
	GetState() (started []string, err error)
	Publish() error
}

type IFxManagerService interface {
	IFxManager
	IFxService.IFxServices
}

type IFxManagerData interface {
	IFxManager
	IDataShutDown.IDataShutDown
}
