package internal

import "github.com/bhbosman/gocommon/model"

type IdAndName struct {
	ServiceId         model.ServiceIdentifier
	ServiceDependency model.ServiceIdentifier
	Name              string
	Active            bool
}

type SendActionsForService struct {
	Name    string
	Actions []string
}
