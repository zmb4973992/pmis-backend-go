package dao

type dao struct {
	departmentDAO
	departmentAndUserDAO
	projectBreakdownDAO
	relatedPartyDAO
	roleAndUserDAO
	userDAO
}

var (
	entranceOfAllDAO     = new(dao)
	DepartmentDAO        = entranceOfAllDAO.departmentDAO
	DepartmentAndUserDAO = entranceOfAllDAO.departmentAndUserDAO
	ProjectBreakdownDAO  = entranceOfAllDAO.projectBreakdownDAO
	RelatedPartyDAO      = entranceOfAllDAO.relatedPartyDAO
	RoleAndUserDAO       = entranceOfAllDAO.roleAndUserDAO
	UserDAO              = entranceOfAllDAO.userDAO
)
