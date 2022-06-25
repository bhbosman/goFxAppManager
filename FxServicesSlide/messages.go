package FxServicesSlide

import "github.com/bhbosman/gocommon/model"

type PublishInstanceDataFor struct {
	Name string
}

type SendActionsForService struct {
	Name    string
	Actions []string
}

type IdAndName struct {
	ServiceId         model.ServiceIdentifier
	ServiceDependency model.ServiceIdentifier
	Name              string
	Active            bool
}

type FxServicesManagerData struct {
	Name              string
	Active            bool
	ServiceId         model.ServiceIdentifier
	ServiceDependency model.ServiceIdentifier
	isDirty           bool
}
