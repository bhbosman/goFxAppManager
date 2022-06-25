all: DeleteAllGeneratedFiles  IFxManagerSlide

DeleteAllGeneratedFiles:
	del *_InterfaceMethods.go




IFxManagerSlide:
	mockgen -package internal -generateWhat ddd -destination IFxManagerSlide_InterfaceMethods.go . IFxManagerSlide
