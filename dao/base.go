package dao

type dao struct {
	departmentDAO
	departmentAndUserDAO
	projectDisassemblyDAO
	relatedPartyDAO
	roleAndUserDAO
	userDAO
	operationRecordDAO
}

var (
	entrance              = new(dao)
	DepartmentDAO         = entrance.departmentDAO
	DepartmentAndUserDAO  = entrance.departmentAndUserDAO
	ProjectDisassemblyDAO = entrance.projectDisassemblyDAO
	RelatedPartyDAO       = entrance.relatedPartyDAO
	RoleAndUserDAO        = entrance.roleAndUserDAO
	UserDAO               = entrance.userDAO
	OperationRecordDAO    = entrance.operationRecordDAO
)
