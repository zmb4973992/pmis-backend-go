package controller

// Base 这里定义controller层的基础方法，避免在其他的controller中反复重写
type Base struct{}

type controller struct {
	captcha
	department
	noRoute
	disassembly
	relatedParty
	user
	operationLog
	roleAndUser
	errorLog
	token
	project
	contract
	dictionaryType
	dictionaryItem
	fileManagement
	progress
}

var (
	entrance        = new(controller)
	Captcha         = entrance.captcha
	NoRoute         = entrance.noRoute
	Department      = entrance.department
	Disassembly     = entrance.disassembly
	RelatedParty    = entrance.relatedParty
	User            = entrance.user
	OperationRecord = entrance.operationLog
	RoleAndUser     = entrance.roleAndUser
	ErrorLog        = entrance.errorLog
	Token           = entrance.token
	Project         = entrance.project
	Contract        = entrance.contract
	DictionaryType  = entrance.dictionaryType
	DictionaryItem  = entrance.dictionaryItem
	FileManagement  = entrance.fileManagement
	Progress        = entrance.progress
)
