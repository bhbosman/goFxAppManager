package internal

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
)

type IFxManagerData interface {
	IFxManagerSlide
	ISendMessage.ISendMessage
	IDataShutDown.IDataShutDown
}
