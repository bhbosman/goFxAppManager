package FxServicesSlide

import "github.com/bhbosman/gocommon/model"

type publishInstanceDataFor struct {
	Name string
}

type FxServicesManagerData struct {
	Name              string
	Active            bool
	ServiceId         model.ServiceIdentifier
	ServiceDependency model.ServiceIdentifier
	isDirty           bool
}
