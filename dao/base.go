package dao

type dao struct {
	departmentDAO
	departmentAndUserDAO
	projectBreakdownDAO
	relatedPartyDAO
	roleAndUserDAO
	userDAO
	operationRecordDAO
}

var (
	entrance             = new(dao)
	DepartmentDAO        = entrance.departmentDAO
	DepartmentAndUserDAO = entrance.departmentAndUserDAO
	ProjectBreakdownDAO  = entrance.projectBreakdownDAO
	RelatedPartyDAO      = entrance.relatedPartyDAO
	RoleAndUserDAO       = entrance.roleAndUserDAO
	UserDAO              = entrance.userDAO
	OperationRecordDAO   = entrance.operationRecordDAO
)
