package model

import "github.com/bhbosman/gocommon/model"

type FxServiceStarted struct {
	Name string
}

type FxServiceStopped struct {
	Name string
}

type FxServiceAdded struct {
	Name string
}

type FxServiceStatus struct {
	Name              string
	Active            bool
	ServiceId         model.ServiceIdentifier
	ServiceDependency model.ServiceIdentifier
}
