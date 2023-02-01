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
	dictionaryItemService
	dictionaryTypeService
}

// 定义各个服务的入口,避免反复new service
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
	DictionaryItemService      = entrance.dictionaryItemService
	DictionaryTypeService      = entrance.dictionaryTypeService
)

var (
	//更新数据库记录时一定需要省略的字段：创建者和删除者的相关字段
	fieldsToBeOmittedWhenUpdating = []string{
		"created_at", "iCreate", "deleted_at", "deleter"}
)
