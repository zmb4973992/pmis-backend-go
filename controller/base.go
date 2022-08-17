package controller

// BaseController 这里定义controller层的基础方法，如success、failure，避免在其他的controller中反复重写
type BaseController struct{}

type controller struct {
	departmentController
	noRouteController
	disassemblyController
	relatedPartyController
	userController
	operationRecordController
	roleAndUserController
}

var (
	entrance                  = new(controller)
	NoRouteController         = entrance.noRouteController
	DepartmentController      = entrance.departmentController
	DisassemblyController     = entrance.disassemblyController
	RelatedPartyController    = entrance.relatedPartyController
	UserController            = entrance.userController
	OperationRecordController = entrance.operationRecordController
	RoleAndUserController     = entrance.roleAndUserController
)
