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
	public.JWTRouter

	private.UserRouter
	private.FileRouter
	private.ContractRouter
	private.RelatedPartyRouter
	private.OrganizationRouter
	private.ProjectRouter
	private.DisassemblyRouter
	private.DictionaryTypeRouter
	private.DictionaryDetailRouter
	private.ProgressRouter
	private.IncomeAndExpenditureRouter
	private.RoleRouter
	private.UserAndRoleRouter
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
	customRouterGroup.InitTestRouter(publicGroup)
	customRouterGroup.InitLoginRouter(publicGroup)
	customRouterGroup.InitStaticRouter(publicGroup)
	customRouterGroup.InitDownloadRouter(publicGroup)
	customRouterGroup.InitRegisterRouter(publicGroup)
	customRouterGroup.InitCaptchaRouter(publicGroup)
	customRouterGroup.InitJWTRouter(publicGroup)

	privateGroup := engine.Group("")
	privateGroup.Use(middleware.RateLimit(), middleware.JWT())
	customRouterGroup.InitUserRouter(privateGroup)
	customRouterGroup.InitFileRouter(privateGroup)
	customRouterGroup.InitRelatedPartyRouter(privateGroup)
	customRouterGroup.InitOrganizationRouter(privateGroup)
	customRouterGroup.InitProjectRouter(privateGroup)
	customRouterGroup.InitContractRouter(privateGroup)
	customRouterGroup.InitDisassemblyRouter(privateGroup)
	customRouterGroup.InitDictionaryTypeRouter(privateGroup)
	customRouterGroup.InitDictionaryDetailRouter(privateGroup)
	customRouterGroup.InitProgressRouter(privateGroup)
	customRouterGroup.InitIncomeAndExpenditureRouter(privateGroup)
	customRouterGroup.InitRoleRouter(privateGroup)
	customRouterGroup.InitUserAndRoleRouter(privateGroup)

	engine.GET("/snow-id", controller.SnowID.Get) //获取雪花id，以后可删
	//依次加载所有的路由组，以下都需要jwt验证
	api := engine.Group("/api")
	api.Use(middleware.RateLimit())
	{
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
