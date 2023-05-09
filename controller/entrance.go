package controller

// Base 这里定义controller层的基础方法，避免在其他的controller中反复重写
type Base struct{}

type controller struct {
	noRoute
	captcha
	organization
	disassembly
	relatedParty
	user
	requestLog
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
	snowID
}

var (
	entrance             = new(controller)
	Captcha              = entrance.captcha
	NoRoute              = entrance.noRoute
	Organization         = entrance.organization
	Disassembly          = entrance.disassembly
	RelatedParty         = entrance.relatedParty
	User                 = entrance.user
	RequestLog           = entrance.requestLog
	ErrorLog             = entrance.errorLog
	Token                = entrance.token
	Project              = entrance.project
	Contract             = entrance.contract
	DictionaryType       = entrance.dictionaryType
	DictionaryDetail     = entrance.dictionaryDetail
	FileManagement       = entrance.fileManagement
	Progress             = entrance.progress
	IncomeAndExpenditure = entrance.incomeAndExpenditure
	SnowID               = entrance.snowID
	Role                 = entrance.role
)
