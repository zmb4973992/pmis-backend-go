package dao

type dao struct {
	departmentDAO
	disassemblyDAO
	relatedPartyDAO
	userDAO
	operationRecordDAO
}

var (
	entrance           = new(dao)
	DepartmentDAO      = entrance.departmentDAO
	DisassemblyDAO     = entrance.disassemblyDAO
	RelatedPartyDAO    = entrance.relatedPartyDAO
	UserDAO            = entrance.userDAO
	OperationRecordDAO = entrance.operationRecordDAO
)
