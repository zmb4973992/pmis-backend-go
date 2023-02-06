package service

// service 所有服务的入口
type service struct {
	roleAndUser
}

// 定义各个服务的入口,避免反复new service
var (
	entrance    = new(service)
	RoleAndUser = entrance.roleAndUser
)
