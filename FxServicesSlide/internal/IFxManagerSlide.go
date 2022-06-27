package internal

type IFxManagerSlide interface {
	StartService(name string)
	StopService(name string)
	StartAllService()
	StopAllService()
	SetConnectionListChange(cb func(connectionList []IdAndName))
	SetConnectionInstanceChange(cb func(data SendActionsForService))
}
