package controller

// Base 这里定义controller层的基础方法，避免在其他的controller中反复重写
type Base struct{}

type controller struct {
	department
	noRoute
	disassembly
	disassemblyTemplate
	relatedParty
	user
	operationRecord
	roleAndUser
	errorLog
	token
	project
	dictionaryType
	dictionaryItem
}

var (
	entrance            = new(controller)
	NoRoute             = entrance.noRoute
	Department          = entrance.department
	Disassembly         = entrance.disassembly
	DisassemblyTemplate = entrance.disassemblyTemplate
	RelatedParty        = entrance.relatedParty
	User                = entrance.user
	OperationRecord     = entrance.operationRecord
	RoleAndUser         = entrance.roleAndUser
	ErrorLog            = entrance.errorLog
	Token               = entrance.token
	Project             = entrance.project
	DictionaryType      = entrance.dictionaryType
	DictionaryItem      = entrance.dictionaryItem
)
