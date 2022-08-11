package router

import (
	"github.com/gin-gonic/gin"
	"learn-go/controller"
	"learn-go/middleware"
)

// Init 初始化路由器,最终返回*gin.Engine类型，给main调用
func Init() *gin.Engine {
	//使用gin框架，生成默认的空引擎
	engine := gin.New()
	//全局中间件
	engine.Use(middleware.ZapLogger(), gin.Recovery(), middleware.Cors())
	engine.POST("/api/login", controller.Login)                //用户登录
	engine.POST("/api/user", controller.UserController.Create) //添加用户
	engine.POST("/upload_single", controller.UploadSingle)     //测试上传单个
	engine.POST("/upload_multiple", controller.UploadMultiple) //测试上传多个

	//依次加载所有的路由组，以下都需要登录验证(jwt验证)
	api := engine.Group("/api").Use(middleware.NeedLogin())
	{
		api.GET("/user/:id", middleware.NeedAuth(), controller.UserController.Get) //获取用户详情
		api.PUT("/user/:id", controller.UserController.Update)                     //修改用户
		api.DELETE("/user/:id", controller.UserController.Delete)                  //删除用户
		api.GET("/user/list", controller.UserController.List)                      //获取用户列表

		api.GET("/related_party/:id", controller.RelatedPartyController.Get)    //获取相关方详情
		api.PUT("/related_party/:id", controller.RelatedPartyController.Update) //修改相关方
		api.POST("/related_party", controller.RelatedPartyController.Create)    //新增相关方
		api.DELETE("/:id", controller.RelatedPartyController.Delete)            //删除相关方
		api.GET("/related_party/list", controller.RelatedPartyController.List)  //获取相关方列表

		api.GET("/department/:id", controller.DepartmentController.Get)       //获取部门详情
		api.POST("/department", controller.DepartmentController.Create)       //新增部门
		api.PUT("/department/:id", controller.DepartmentController.Update)    //修改部门
		api.DELETE("/department/:id", controller.DepartmentController.Delete) //删除部门
		api.GET("/department/list", controller.DepartmentController.List)     //获取部门列表

		api.GET("/project_breakdown/:id", controller.ProjectBreakdownController.Get)       //获取项目拆解详情
		api.POST("/project_breakdown", controller.ProjectBreakdownController.Create)       //新增项目拆解
		api.PUT("/project_breakdown/:id", controller.ProjectBreakdownController.Update)    //修改项目拆解
		api.DELETE("/project_breakdown/:id", controller.ProjectBreakdownController.Delete) //删除项目拆解
		api.GET("/project_breakdown/list", controller.ProjectBreakdownController.List)     //获取项目拆解列表
	}

	engine.NoRoute(controller.NoRouteController.NoRoute)

	//引擎处理完成后，返回
	return engine
}
