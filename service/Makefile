all: DeleteAllGeneratedFiles IFxManager

DeleteAllGeneratedFiles:
	del *_InterfaceMethods.go



IFxManager:
	mockgen -package service -generateWhat ddd -destination IFxManager_InterfaceMethods.go . IFxManager
	mockgen -package service -generateWhat mockgen -destination IFxManagerMock.go . IFxManager

