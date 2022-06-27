package internal

import (
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
)

type IFxManagerService interface {
	IFxManagerSlide
	ISendMessage.ISendMessage
	IFxService.IFxServices
}
