all: DeleteAllGeneratedFiles IFxManager

DeleteAllGeneratedFiles:
	del *_InterfaceMethods.go



IFxManager:
	mockgen -package FxAppManagerService -generateWhat ddd -destination IFxManager_InterfaceMethods.go . IFxManager
	mockgen -package FxAppManagerService -generateWhat mockgen -destination IFxManagerMock.go . IFxManager

