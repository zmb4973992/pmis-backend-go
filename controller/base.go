package controller

// BaseController 这里定义controller层的基础方法，如success、failure，避免在其他的controller中反复重写
type BaseController struct{}

type controller struct {
	departmentController
	noRouteController
	disassemblyController
	disassemblyTemplateController
	relatedPartyController
	userController
	operationRecordController
	roleAndUserController
	errorLogController
	tokenController
	projectController
}

var (
	entrance                      = new(controller)
	NoRouteController             = entrance.noRouteController
	DepartmentController          = entrance.departmentController
	DisassemblyController         = entrance.disassemblyController
	DisassemblyTemplateController = entrance.disassemblyTemplateController
	RelatedPartyController        = entrance.relatedPartyController
	UserController                = entrance.userController
	OperationRecordController     = entrance.operationRecordController
	RoleAndUserController         = entrance.roleAndUserController
	ErrorLogController            = entrance.errorLogController
	TokenController               = entrance.tokenController
	ProjectController             = entrance.projectController
)
