package dao

type dao struct {
	departmentDAO
	projectDisassemblyDAO
	relatedPartyDAO
	userDAO
	operationRecordDAO
}

var (
	entrance              = new(dao)
	DepartmentDAO         = entrance.departmentDAO
	ProjectDisassemblyDAO = entrance.projectDisassemblyDAO
	RelatedPartyDAO       = entrance.relatedPartyDAO
	UserDAO               = entrance.userDAO
	OperationRecordDAO    = entrance.operationRecordDAO
)
