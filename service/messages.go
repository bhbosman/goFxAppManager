package service

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
	Name   string
	Active bool
}
