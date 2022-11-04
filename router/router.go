package router

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
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

	engine.GET("/api/validate_token/:token", controller.TokenController.Validate) //单独校验token是否有效

	//依次加载所有的路由组，以下都需要登录验证(jwt验证)
	api := engine.Group("/api").Use(middleware.NeedLogin())
	{
		api.GET("/user/:id", middleware.NeedAuth(), controller.UserController.Get) //获取用户详情
		api.PUT("/user/:id", controller.UserController.Update)                     //修改用户（目前为全功能，考虑改成：修改用户基本信息）
		api.DELETE("/user/:id", controller.UserController.Delete)                  //删除用户
		api.GET("/user/list", controller.UserController.List)                      //获取用户列表

		api.GET("/user", controller.UserController.GetByToken) //根据header里的token获取用户详情

		api.GET("/role_and_user/role/:role_id", controller.RoleAndUserController.ListByRoleID)      //根据角色ID获取角色和用户的列表
		api.POST("/role_and_user/role/:role_id", controller.RoleAndUserController.CreateByRoleID)   //根据角色ID批量新增角色和用户
		api.PUT("/role_and_user/role/:role_id", controller.RoleAndUserController.UpdateByRoleID)    //根据角色ID批量修改角色和用户
		api.DELETE("/role_and_user/role/:role_id", controller.RoleAndUserController.DeleteByRoleID) //根据角色ID批量删除角色和用户

		api.GET("/role_and_user/user/:user_id", controller.RoleAndUserController.ListByUserID)      //根据用户ID获取角色和用户的列表
		api.POST("/role_and_user/user/:user_id", controller.RoleAndUserController.CreateByUserID)   //根据用户ID批量新增角色和用户
		api.PUT("/role_and_user/user/:user_id", controller.RoleAndUserController.UpdateByUserID)    //根据用户ID批量修改角色和用户
		api.DELETE("/role_and_user/user/:user_id", controller.RoleAndUserController.DeleteByUserID) //根据用户ID批量删除角色和用户

		api.GET("/role_and_user/by_token_in_header", controller.RoleAndUserController.ListByTokenInHeader) //根据header里的token获取角色和用户的列表

		api.GET("/related_party/:id", controller.RelatedPartyController.Get)       //获取相关方详情
		api.PUT("/related_party/:id", controller.RelatedPartyController.Update)    //修改相关方
		api.POST("/related_party", controller.RelatedPartyController.Create)       //新增相关方
		api.DELETE("/related_party/:id", controller.RelatedPartyController.Delete) //删除相关方
		api.GET("/related_party/list", controller.RelatedPartyController.List)     //获取相关方列表

		api.GET("/department/:id", controller.DepartmentController.Get)       //获取部门详情
		api.POST("/department", controller.DepartmentController.Create)       //新增部门
		api.PUT("/department/:id", controller.DepartmentController.Update)    //修改部门
		api.DELETE("/department/:id", controller.DepartmentController.Delete) //删除部门
		api.GET("/department/list", controller.DepartmentController.List)     //获取部门列表

		api.GET("/disassembly/:id", controller.DisassemblyController.Get)                //获取项目拆解详情
		api.POST("/disassembly", controller.DisassemblyController.Create)                //新增项目拆解
		api.POST("/disassembly/batch", controller.DisassemblyController.CreateInBatches) //批量新增项目拆解
		api.PUT("/disassembly/:id", controller.DisassemblyController.Update)             //修改项目拆解
		api.DELETE("/disassembly/:id", controller.DisassemblyController.Delete)          //删除项目拆解
		api.GET("/disassembly/list", controller.DisassemblyController.List)              //获取项目拆解列表

		api.GET("/disassembly_template/:id", controller.DisassemblyTemplateController.Get)       //获取项目拆解模板详情
		api.POST("/disassembly_template", controller.DisassemblyTemplateController.Create)       //新增项目拆解模板
		api.PUT("/disassembly_template/:id", controller.DisassemblyTemplateController.Update)    //修改项目拆解模板
		api.DELETE("/disassembly_template/:id", controller.DisassemblyTemplateController.Delete) //删除项目拆解模板
		api.GET("/disassembly_template/list", controller.DisassemblyTemplateController.List)     //获取项目拆解模板列表

		api.GET("/operation_record/:id", controller.OperationRecordController.Get)       //获取操作记录详情
		api.POST("/operation_record", controller.OperationRecordController.Create)       //新增操作记录
		api.PUT("/operation_record/:id", controller.OperationRecordController.Update)    //修改操作记录
		api.DELETE("/operation_record/:id", controller.OperationRecordController.Delete) //删除操作记录
		api.GET("/operation_record/list", controller.OperationRecordController.List)     //获取操作详情列表

		api.GET("/error_log/:id", controller.ErrorLogController.Get)       //获取错误日志详情
		api.POST("/error_log", controller.ErrorLogController.Create)       //新增错误日志
		api.PUT("/error_log/:id", controller.ErrorLogController.Update)    //修改错误日志
		api.DELETE("/error_log/:id", controller.ErrorLogController.Delete) //删除错误日志
	}

	engine.NoRoute(controller.NoRouteController.NoRoute)

	//引擎处理完成后，返回
	return engine
}
