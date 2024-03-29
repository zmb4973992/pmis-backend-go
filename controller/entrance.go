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
	file
	progress
	incomeAndExpenditure
	role
	menu
	projectDailyAndCumulativeIncome
	projectDailyAndCumulativeExpenditure
	contractDailyAndCumulativeIncome
	contractDailyAndCumulativeExpenditure
	message
	operationLog
}

var (
	entrance                              = new(controller)
	Captcha                               = entrance.captcha
	NoRoute                               = entrance.noRoute
	Organization                          = entrance.organization
	Disassembly                           = entrance.disassembly
	RelatedParty                          = entrance.relatedParty
	User                                  = entrance.user
	RequestLog                            = entrance.requestLog
	ErrorLog                              = entrance.errorLog
	Token                                 = entrance.token
	Project                               = entrance.project
	Contract                              = entrance.contract
	DictionaryType                        = entrance.dictionaryType
	DictionaryDetail                      = entrance.dictionaryDetail
	File                                  = entrance.file
	Progress                              = entrance.progress
	IncomeAndExpenditure                  = entrance.incomeAndExpenditure
	Role                                  = entrance.role
	Menu                                  = entrance.menu
	ProjectDailyAndCumulativeIncome       = entrance.projectDailyAndCumulativeIncome
	ProjectDailyAndCumulativeExpenditure  = entrance.projectDailyAndCumulativeExpenditure
	ContractDailyAndCumulativeIncome      = entrance.contractDailyAndCumulativeIncome
	ContractDailyAndCumulativeExpenditure = entrance.contractDailyAndCumulativeExpenditure
	Message                               = entrance.message
	OperationLog                          = entrance.operationLog
)
