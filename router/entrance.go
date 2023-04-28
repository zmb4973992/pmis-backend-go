package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/global"
	"pmis-backend-go/middleware"
	"pmis-backend-go/router/private"
	"pmis-backend-go/router/public"
)

type CustomRouterGroup struct {
	public.LoginRouter
	public.StaticRouter
	public.TestRouter
	public.DownloadRouter
	public.RegisterRouter
	public.CaptchaRouter
	public.ValidateTokenRouter

	private.UserRouter
	private.FileRouter
	private.ContractRouter
	private.RelatedPartyRouter
	private.DepartmentRouter
	private.ProjectRouter
	private.DisassemblyRouter
	private.DictionaryTypeRouter
	private.DictionaryDetailRouter
	private.ProgressRouter
	private.IncomeAndExpenditureRouter
	private.RoleRouter
}

// InitEngine 初始化路由器,最终返回*gin.Engine类型，给main调用
func InitEngine() *gin.Engine {
	//设置运行模式
	gin.SetMode(global.Config.AppConfig.AppMode)
	fmt.Println("当前运行模式为：", gin.Mode())
	engine := gin.New()

	//全局中间件
	engine.Use(middleware.Logger(), gin.Recovery(), middleware.Cors())

	//路由不匹配时的处理
	engine.NoRoute(controller.NoRoute.NoRoute)

	//将待处理的路由分为两类：不需要jwt鉴权的为public，需要jwt鉴权的为private
	publicGroup := engine.Group("")
	customRouterGroup := new(CustomRouterGroup)
	customRouterGroup.InitLoginRouter(publicGroup)
	customRouterGroup.InitStaticRouter(publicGroup)
	customRouterGroup.InitDownloadRouter(publicGroup)
	customRouterGroup.InitRegisterRouter(publicGroup)
	customRouterGroup.InitCaptchaRouter(publicGroup)
	customRouterGroup.InitValidateTokenRouter(publicGroup)

	privateGroup := engine.Group("")
	privateGroup.Use(middleware.ValidateToken(), middleware.RateLimit())
	customRouterGroup.InitUserRouter(privateGroup)
	customRouterGroup.InitFileRouter(privateGroup)
	customRouterGroup.InitRelatedPartyRouter(privateGroup)
	customRouterGroup.InitDepartmentRouter(privateGroup)
	customRouterGroup.InitProjectRouter(privateGroup)
	customRouterGroup.InitContractRouter(privateGroup)
	customRouterGroup.InitDisassemblyRouter(privateGroup)
	customRouterGroup.InitDictionaryTypeRouter(privateGroup)
	customRouterGroup.InitDictionaryDetailRouter(privateGroup)
	customRouterGroup.InitProgressRouter(privateGroup)
	customRouterGroup.InitIncomeAndExpenditureRouter(privateGroup)
	customRouterGroup.InitRoleRouter(privateGroup)

	//依次加载所有的路由组，以下都需要jwt验证
	api := engine.Group("/api")
	api.Use(middleware.ValidateToken(), middleware.RateLimit())
	{
		roleAndUser := api.Group("/role-and-user")
		{
			roleAndUser.GET("/role/:role-id", controller.RoleAndUser.ListByRoleID)      //根据角色ID获取角色和用户的列表
			roleAndUser.POST("/role/:role-id", controller.RoleAndUser.CreateByRoleID)   //根据角色ID批量新增角色和用户
			roleAndUser.PUT("/role/:role-id", controller.RoleAndUser.UpdateByRoleID)    //根据角色ID批量修改角色和用户
			roleAndUser.DELETE("/role/:role-id", controller.RoleAndUser.DeleteByRoleID) //根据角色ID批量删除角色和用户

			roleAndUser.GET("/user/:user-id", controller.RoleAndUser.ListByUserID)      //根据用户ID获取角色和用户的列表
			roleAndUser.POST("/user/:user-id", controller.RoleAndUser.CreateByUserID)   //根据用户ID批量新增角色和用户
			roleAndUser.PUT("/user/:user-id", controller.RoleAndUser.UpdateByUserID)    //根据用户ID批量修改角色和用户
			roleAndUser.DELETE("/user/:user-id", controller.RoleAndUser.DeleteByUserID) //根据用户ID批量删除角色和用户

			roleAndUser.GET("/by-token-in-header", controller.RoleAndUser.ListByTokenInHeader) //根据header里的token获取角色和用户的列表
		}
		operationLog := api.Group("/operation-log")
		{
			operationLog.GET("/:operation-log-id", controller.OperationRecord.Get)       //获取操作记录详情
			operationLog.DELETE("/:operation-log-id", controller.OperationRecord.Delete) //删除操作记录
			operationLog.POST("/list", controller.OperationRecord.GetList)               //获取操作详情列表
		}
		errorLog := api.Group("/error-log")
		{
			errorLog.GET("/:error-log-id", controller.ErrorLog.Get)       //获取错误日志详情
			errorLog.POST("", controller.ErrorLog.Create)                 //新增错误日志
			errorLog.PATCH("/:error-log-id", controller.ErrorLog.Update)  //修改错误日志
			errorLog.DELETE("/:error-log-id", controller.ErrorLog.Delete) //删除错误日志
		}
	}

	//引擎配置完成后，返回
	return engine
}
