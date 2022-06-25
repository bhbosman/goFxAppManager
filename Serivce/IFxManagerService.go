package Serivce

import (
	"github.com/bhbosman/goFxAppManager/Serivce/internal"
	"github.com/bhbosman/gocommon/Services/IFxService"
)

type IFxManagerService interface {
	internal.IFxManager
	IFxService.IFxServices
}
