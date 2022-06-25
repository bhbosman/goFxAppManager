package internal

import "github.com/bhbosman/gocommon/Services/ISendMessage"

type IFxManagerSlide interface {
	StartService(name string)
	StopService(name string)
	StartAllService()
	StopAllService()
}

type IFxManagerData interface {
	IFxManagerSlide
	ISendMessage.ISendMessage
}
