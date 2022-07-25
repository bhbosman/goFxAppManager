package internal

type IdAndName struct {
	Name   string
	Active bool
}

type SendActionsForService struct {
	Name    string
	Actions []string
}
