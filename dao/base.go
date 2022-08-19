package dao

type dao struct {
	departmentDAO
	disassemblyDAO
	relatedPartyDAO
	userDAO
	operationRecordDAO
	disassemblyTemplateDAO
}

var (
	entrance               = new(dao)
	DepartmentDAO          = entrance.departmentDAO
	DisassemblyDAO         = entrance.disassemblyDAO
	DisassemblyTemplateDAO = entrance.disassemblyTemplateDAO
	RelatedPartyDAO        = entrance.relatedPartyDAO
	UserDAO                = entrance.userDAO
	OperationRecordDAO     = entrance.operationRecordDAO
)
