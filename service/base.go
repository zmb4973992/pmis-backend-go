package service

type baseService struct{}

// EntranceOfAllService 所有服务的入口
type EntranceOfAllService struct {
	loginService
	userService
	relatedPartyService
	departmentService
	projectBreakdownService
}

//定义各个服务的入口,避免反复new service
var (
	entrance                = new(EntranceOfAllService)
	LoginService            = entrance.loginService
	UserService             = entrance.userService
	RelatedPartyService     = entrance.relatedPartyService
	DepartmentService       = entrance.departmentService
	ProjectBreakdownService = entrance.projectBreakdownService
)
