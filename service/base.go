package service

// AllService 所有服务的入口
type AllService struct {
	loginService
	userService
	relatedPartyService
	departmentService
	disassemblyService
	operationRecordService
	roleAndUserService
	disassemblyTemplateService
	errorLogService
	projectService
}

//定义各个服务的入口,避免反复new service
var (
	entrance                   = new(AllService)
	LoginService               = entrance.loginService
	UserService                = entrance.userService
	RelatedPartyService        = entrance.relatedPartyService
	DepartmentService          = entrance.departmentService
	DisassemblyService         = entrance.disassemblyService
	DisassemblyTemplateService = entrance.disassemblyTemplateService
	OperationRecordService     = entrance.operationRecordService
	RoleAndUserService         = entrance.roleAndUserService
	ErrorLogService            = entrance.errorLogService
	ProjectService             = entrance.projectService
)
