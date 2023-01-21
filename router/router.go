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
	engine.POST("/upload-single", controller.UploadSingle)     //测试上传单个
	engine.POST("/upload-multiple", controller.UploadMultiple) //测试上传多个

	engine.GET("/api/validate-token/:token", controller.TokenController.Validate) //单独校验token是否有效

	//依次加载所有的路由组，以下都需要登录验证(jwt验证)
	api := engine.Group("/api").Use(middleware.NeedLogin(), middleware.RateLimit())
	{
		api.GET("/user/:user-id", middleware.NeedAuth(), controller.UserController.Get) //获取用户详情
		api.PUT("/user/:user-id", controller.UserController.Update)                     //修改用户（目前为全功能，考虑改成：修改用户基本信息）
		api.DELETE("/user/:user-id", controller.UserController.Delete)                  //删除用户
		api.POST("/user/list", controller.UserController.List)                          //获取用户列表

		api.GET("/user", controller.UserController.GetByToken) //根据header里的token获取用户详情

		api.GET("/role-and-user/role/:role-id", controller.RoleAndUserController.ListByRoleID)      //根据角色ID获取角色和用户的列表
		api.POST("/role-and-user/role/:role-id", controller.RoleAndUserController.CreateByRoleID)   //根据角色ID批量新增角色和用户
		api.PUT("/role-and-user/role/:role-id", controller.RoleAndUserController.UpdateByRoleID)    //根据角色ID批量修改角色和用户
		api.DELETE("/role-and-user/role/:role-id", controller.RoleAndUserController.DeleteByRoleID) //根据角色ID批量删除角色和用户

		api.GET("/role-and-user/user/:user-id", controller.RoleAndUserController.ListByUserID)      //根据用户ID获取角色和用户的列表
		api.POST("/role-and-user/user/:user-id", controller.RoleAndUserController.CreateByUserID)   //根据用户ID批量新增角色和用户
		api.PUT("/role-and-user/user/:user-id", controller.RoleAndUserController.UpdateByUserID)    //根据用户ID批量修改角色和用户
		api.DELETE("/role-and-user/user/:user-id", controller.RoleAndUserController.DeleteByUserID) //根据用户ID批量删除角色和用户

		api.GET("/role-and-user/by-token-in-header", controller.RoleAndUserController.ListByTokenInHeader) //根据header里的token获取角色和用户的列表

		api.GET("/related-party/:related-party-id", controller.RelatedPartyController.Get)       //获取相关方详情
		api.PUT("/related-party/:related-party-id", controller.RelatedPartyController.Update)    //修改相关方
		api.POST("/related-party", controller.RelatedPartyController.Create)                     //新增相关方
		api.DELETE("/related-party/:related-party-id", controller.RelatedPartyController.Delete) //删除相关方
		api.POST("/related-party/list", controller.RelatedPartyController.List)                  //获取相关方列表

		api.GET("/department/:department-id", controller.DepartmentController.Get)       //获取部门详情
		api.POST("/department", controller.DepartmentController.Create)                  //新增部门
		api.PUT("/department/:department-id", controller.DepartmentController.Update)    //修改部门
		api.DELETE("/department/:department-id", controller.DepartmentController.Delete) //删除部门
		api.POST("/department/list", controller.DepartmentController.List)               //获取部门列表

		api.GET("/project/:project-id", controller.ProjectController.Get)        //获取项目详情
		api.POST("/project", controller.ProjectController.Create)                //新增项目
		api.POST("/project/batch", controller.ProjectController.CreateInBatches) //批量新增项目
		api.PUT("/project/:project-id", controller.ProjectController.Update)     //修改项目
		api.DELETE("/project/:project-id", controller.ProjectController.Delete)  //删除项目
		api.POST("project/list", controller.ProjectController.List)              //获取项目列表

		api.GET("/disassembly/:disassembly-id", controller.DisassemblyController.Get)                           //获取项目拆解详情
		api.POST("/disassembly/tree", controller.DisassemblyController.Tree)                                    //获取项目拆解的节点树
		api.POST("/disassembly", controller.DisassemblyController.Create)                                       //新增项目拆解
		api.POST("/disassembly/batch", controller.DisassemblyController.CreateInBatches)                        //批量新增项目拆解
		api.PUT("/disassembly/:disassembly-id", controller.DisassemblyController.Update)                        //修改项目拆解
		api.DELETE("/disassembly/:disassembly-id", controller.DisassemblyController.Delete)                     //删除项目拆解
		api.DELETE("/disassembly/cascade/:disassembly-id", controller.DisassemblyController.DeleteWithSubitems) //删除项目拆解（子项一并删除）
		api.POST("/disassembly/list", controller.DisassemblyController.List)                                    //获取项目拆解列表

		api.GET("/disassembly-template/:disassembly-template-id", controller.DisassemblyTemplateController.Get)       //获取项目拆解模板详情
		api.POST("/disassembly-template", controller.DisassemblyTemplateController.Create)                            //新增项目拆解模板
		api.PUT("/disassembly-template/:disassembly-template-id", controller.DisassemblyTemplateController.Update)    //修改项目拆解模板
		api.DELETE("/disassembly-template/:disassembly-template-id", controller.DisassemblyTemplateController.Delete) //删除项目拆解模板
		api.POST("/disassembly-template/list", controller.DisassemblyTemplateController.List)                         //获取项目拆解模板列表

		api.GET("/operation-record/:operation-record-id", controller.OperationRecordController.Get)       //获取操作记录详情
		api.POST("/operation-record", controller.OperationRecordController.Create)                        //新增操作记录
		api.PUT("/operation-record/:operation-record-id", controller.OperationRecordController.Update)    //修改操作记录
		api.DELETE("/operation-record/:operation-record-id", controller.OperationRecordController.Delete) //删除操作记录
		api.POST("/operation-record/list", controller.OperationRecordController.List)                     //获取操作详情列表

		api.GET("/error-log/:error-log-id", controller.ErrorLogController.Get)       //获取错误日志详情
		api.POST("/error-log", controller.ErrorLogController.Create)                 //新增错误日志
		api.PUT("/error-log/:error-log-id", controller.ErrorLogController.Update)    //修改错误日志
		api.DELETE("/error-log/:error-log-id", controller.ErrorLogController.Delete) //删除错误日志

		//数据字典的类型
		api.POST("/dictionary-type", controller.DictionaryTypeController.Create)                       //新增字典类型
		api.POST("/dictionary-type/batch", controller.DictionaryTypeController.CreateInBatches)        //批量新增字典类型
		api.PUT("/dictionary-type/:dictionary-type-id", controller.DictionaryTypeController.Update)    //修改字典类型
		api.DELETE("/dictionary-type/:dictionary-type-id", controller.DictionaryTypeController.Delete) //删除字典类型
		api.POST("/dictionary-type/list", controller.DictionaryTypeController.List)                    //获取字典类型的列表

		//数据字典的详情项
		api.GET("/dictionary-item/:dictionary-type-id", controller.DictionaryItemController.Get)       //获取单个字典项的所有值
		api.POST("/dictionary-item", controller.DictionaryItemController.Create)                       //新增字典项的单个值
		api.POST("/dictionary-item/batch", controller.DictionaryItemController.CreateInBatches)        //新增字典项的多个值
		api.PUT("/dictionary-item/:dictionary-item-id", controller.DictionaryItemController.Update)    //修改字典项的单个值
		api.DELETE("/dictionary-item/:dictionary-item-id", controller.DictionaryItemController.Delete) //删除字典项的单个值
	}

	engine.NoRoute(controller.NoRouteController.NoRoute)

	//引擎处理完成后，返回
	return engine
}
