package controller

// BaseController 这里定义controller层的基础方法，如success、failure，避免在其他的controller中反复重写
type BaseController struct{}

type controller struct {
	departmentController
	noRouteController
	projectBreakdownController
	relatedPartyController
	userController
}

var (
	entranceOfAllControllers   = new(controller)
	NoRouteController          = entranceOfAllControllers.noRouteController
	DepartmentController       = entranceOfAllControllers.departmentController
	ProjectBreakdownController = entranceOfAllControllers.projectBreakdownController
	RelatedPartyController     = entranceOfAllControllers.relatedPartyController
	UserController             = entranceOfAllControllers.userController
)
