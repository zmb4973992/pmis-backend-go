package service

// service 所有服务的入口
type service struct {
	login
	user
	relatedParty
	department
	disassembly
	operationRecord
	roleAndUser
	errorLog
	project
	dictionaryItem
	dictionaryType
}

// 定义各个服务的入口,避免反复new service
var (
	entrance        = new(service)
	Login           = entrance.login
	User            = entrance.user
	RelatedParty    = entrance.relatedParty
	Department      = entrance.department
	Disassembly     = entrance.disassembly
	OperationRecord = entrance.operationRecord
	RoleAndUser     = entrance.roleAndUser
	ErrorLog        = entrance.errorLog
	Project         = entrance.project
	DictionaryItem  = entrance.dictionaryItem
	DictionaryType  = entrance.dictionaryType
)

var (
	//更新数据库记录时一定需要省略的字段：创建者和删除者的相关字段
	fieldsToBeOmittedWhenUpdating = []string{
		"created_at", "iCreate", "deleted_at", "deleter"}
)
