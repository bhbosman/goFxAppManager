package internal

import (
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type IFxManagerService interface {
	IFxManagerSlide
	ISendMessage.ISendMessage
	IFxService.IFxServices
}
