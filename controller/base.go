package controller

// Base 这里定义controller层的基础方法，避免在其他的controller中反复重写
type Base struct{}

type controller struct {
	captcha
	organization
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
	dictionaryDetail
	fileManagement
	progress
	incomeAndExpenditure
	role
}

var (
	entrance             = new(controller)
	Captcha              = entrance.captcha
	NoRoute              = entrance.noRoute
	Organization         = entrance.organization
	Disassembly          = entrance.disassembly
	RelatedParty         = entrance.relatedParty
	User                 = entrance.user
	OperationRecord      = entrance.operationLog
	RoleAndUser          = entrance.roleAndUser
	ErrorLog             = entrance.errorLog
	Token                = entrance.token
	Project              = entrance.project
	Contract             = entrance.contract
	DictionaryType       = entrance.dictionaryType
	DictionaryDetail     = entrance.dictionaryDetail
	FileManagement       = entrance.fileManagement
	Progress             = entrance.progress
	IncomeAndExpenditure = entrance.incomeAndExpenditure
)
