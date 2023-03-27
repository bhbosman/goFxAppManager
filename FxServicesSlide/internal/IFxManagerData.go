package internal

import (
	"github.com/bhbosman/gocommon/services/IDataShutDown"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type IFxManagerData interface {
	IFxManagerSlide
	ISendMessage.ISendMessage
	IDataShutDown.IDataShutDown
}
