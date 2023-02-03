package service

// service 所有服务的入口
type service struct {
	login
	user
	relatedParty
	department
	disassembly
	operationLog
	roleAndUser
	errorLog
	project
	dictionaryItem
}

// 定义各个服务的入口,避免反复new service
var (
	entrance       = new(service)
	Login          = entrance.login
	User           = entrance.user
	RelatedParty   = entrance.relatedParty
	Department     = entrance.department
	Disassembly    = entrance.disassembly
	OperationLog   = entrance.operationLog
	RoleAndUser    = entrance.roleAndUser
	ErrorLog       = entrance.errorLog
	Project        = entrance.project
	DictionaryItem = entrance.dictionaryItem
)
