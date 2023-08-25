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
	public.RegisterRouter
	public.CaptchaRouter
	public.JWTRouter
	public.DownLoadRouter

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
	private.RequestLogRouter
	private.ErrorLogRouter
	private.MenuRouter
	private.ProjectDailyAndCumulativeIncomeRouter
	private.ProjectDailyAndCumulativeExpenditureRouter
	private.ContractDailyAndCumulativeIncomeRouter
	private.ContractDailyAndCumulativeExpenditureRouter
	private.MessageRouter
}

// InitEngine 初始化路由器,最终返回*gin.Engine类型，给main调用
func InitEngine() *gin.Engine {
	//设置运行模式
	gin.SetMode(global.Config.AppConfig.AppMode)
	fmt.Println("当前运行模式为：", gin.Mode())
	engine := gin.New()

	//全局中间件
	engine.Use(
		gin.Recovery(),
		middleware.Logger(),
		middleware.Cors(),
		middleware.RequestLog(),
		middleware.RateLimit(),
	)

	//路由不匹配时的处理
	engine.NoRoute(controller.NoRoute.NoRoute)

	//将待处理的路由分为两类：不需要jwt鉴权的为public，需要jwt鉴权的为private
	publicGroup := engine.Group("")
	customRouterGroup := new(CustomRouterGroup)
	customRouterGroup.InitTestRouter(publicGroup)
	customRouterGroup.InitLoginRouter(publicGroup)
	customRouterGroup.InitStaticRouter(publicGroup)
	customRouterGroup.InitRegisterRouter(publicGroup)
	customRouterGroup.InitCaptchaRouter(publicGroup)
	customRouterGroup.InitJWTRouter(publicGroup)
	customRouterGroup.InitDownLoadRouter(publicGroup)

	privateGroup := engine.Group("")
	privateGroup.Use(middleware.JWT())
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
	customRouterGroup.InitRequestLogRouter(privateGroup)
	customRouterGroup.InitErrorLogRouter(privateGroup)
	customRouterGroup.InitMenuRouter(privateGroup)
	customRouterGroup.InitProjectDailyAndCumulativeIncomeRouter(privateGroup)
	customRouterGroup.InitProjectDailyAndCumulativeExpenditureRouter(privateGroup)
	customRouterGroup.InitContractDailyAndCumulativeIncomeRouter(privateGroup)
	customRouterGroup.InitContractDailyAndCumulativeExpenditureRouter(privateGroup)
	customRouterGroup.InitMessageRouter(privateGroup)

	//引擎配置完成后，返回
	return engine
}
