package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/global"
	"pmis-backend-go/middleware"
)

// Init 初始化路由器,最终返回*gin.Engine类型，给main调用
func Init() *gin.Engine {
	//设置运行模式
	gin.SetMode(global.Config.APPConfig.AppMode)
	fmt.Println("当前运行模式为：", gin.Mode())

	engine := gin.New()

	//将目录(root变量)下的所有文件设置为静态文件，可以直接访问
	engine.Static("/s", "./static")

	//路由不匹配时的处理
	engine.NoRoute(controller.NoRouteController.NoRoute)

	//全局中间件
	engine.Use(middleware.ZapLogger(), gin.Recovery(), middleware.Cors())

	engine.POST("/api/login", controller.Login)                //用户登录
	engine.POST("/api/user", controller.UserController.Create) //添加用户
	engine.POST("/upload-single", controller.UploadSingle)     //测试上传单个
	engine.POST("/upload-multiple", controller.UploadMultiple) //测试上传多个

	engine.GET("/api/validate-token/:token", controller.TokenController.Validate) //单独校验token是否有效

	//依次加载所有的路由组，以下都需要jwt验证
	api := engine.Group("/api")
	api.Use(middleware.NeedLogin(), middleware.RateLimit())
	{
		user := api.Group("/user")
		{
			user.GET("/:user-id", middleware.NeedAuth(), controller.UserController.Get) //获取用户详情
			user.GET("", controller.UserController.GetByToken)                          //根据header里的token获取用户详情
			user.PUT("/:user-id", controller.UserController.Update)                     //修改用户（目前为全功能，考虑改成：修改用户基本信息）
			user.DELETE("/:user-id", controller.UserController.Delete)                  //删除用户
			user.POST("/list", controller.UserController.List)                          //获取用户列表
		}
		roleAndUser := api.Group("/role-and-user")
		{
			roleAndUser.GET("/role/:role-id", controller.RoleAndUserController.ListByRoleID)      //根据角色ID获取角色和用户的列表
			roleAndUser.POST("/role/:role-id", controller.RoleAndUserController.CreateByRoleID)   //根据角色ID批量新增角色和用户
			roleAndUser.PUT("/role/:role-id", controller.RoleAndUserController.UpdateByRoleID)    //根据角色ID批量修改角色和用户
			roleAndUser.DELETE("/role/:role-id", controller.RoleAndUserController.DeleteByRoleID) //根据角色ID批量删除角色和用户

			roleAndUser.GET("/user/:user-id", controller.RoleAndUserController.ListByUserID)      //根据用户ID获取角色和用户的列表
			roleAndUser.POST("/user/:user-id", controller.RoleAndUserController.CreateByUserID)   //根据用户ID批量新增角色和用户
			roleAndUser.PUT("/user/:user-id", controller.RoleAndUserController.UpdateByUserID)    //根据用户ID批量修改角色和用户
			roleAndUser.DELETE("/user/:user-id", controller.RoleAndUserController.DeleteByUserID) //根据用户ID批量删除角色和用户

			roleAndUser.GET("/by-token-in-header", controller.RoleAndUserController.ListByTokenInHeader) //根据header里的token获取角色和用户的列表
		}
		relatedParty := api.Group("/related-party")
		{
			relatedParty.GET("/:related-party-id", controller.RelatedPartyController.Get)       //获取相关方详情
			relatedParty.PUT("/:related-party-id", controller.RelatedPartyController.Update)    //修改相关方
			relatedParty.POST("", controller.RelatedPartyController.Create)                     //新增相关方
			relatedParty.DELETE("/:related-party-id", controller.RelatedPartyController.Delete) //删除相关方
			relatedParty.POST("/list", controller.RelatedPartyController.List)                  //获取相关方列表
		}
		department := api.Group("/department")
		{
			department.GET("/:department-id", controller.DepartmentController.Get)       //获取部门详情
			department.POST("", controller.DepartmentController.Create)                  //新增部门
			department.PUT("/:department-id", controller.DepartmentController.Update)    //修改部门
			department.DELETE("/:department-id", controller.DepartmentController.Delete) //删除部门
			department.POST("/list", controller.DepartmentController.List)               //获取部门列表
		}
		project := api.Group("/project")
		{
			project.GET("/:project-id", controller.ProjectController.Get)        //获取项目详情
			project.POST("", controller.ProjectController.Create)                //新增项目
			project.POST("/batch", controller.ProjectController.CreateInBatches) //批量新增项目
			project.PUT("/:project-id", controller.ProjectController.Update)     //修改项目
			project.DELETE("/:project-id", controller.ProjectController.Delete)  //删除项目
			project.POST("/list", controller.ProjectController.List)             //获取项目列表
		}
		disassembly := api.Group("/disassembly")
		{
			disassembly.GET("/:disassembly-id", controller.DisassemblyController.Get)                           //获取项目拆解详情
			disassembly.POST("/tree", controller.DisassemblyController.Tree)                                    //获取项目拆解的节点树
			disassembly.POST("", controller.DisassemblyController.Create)                                       //新增项目拆解
			disassembly.POST("/batch", controller.DisassemblyController.CreateInBatches)                        //批量新增项目拆解
			disassembly.PUT("/:disassembly-id", controller.DisassemblyController.Update)                        //修改项目拆解
			disassembly.DELETE("/:disassembly-id", controller.DisassemblyController.Delete)                     //删除项目拆解
			disassembly.DELETE("/cascade/:disassembly-id", controller.DisassemblyController.DeleteWithSubitems) //删除项目拆解（子项一并删除）
			disassembly.POST("/list", controller.DisassemblyController.List)                                    //获取项目拆解列表
		}
		disassemblyTemplate := api.Group("/disassembly-template")
		{
			disassemblyTemplate.GET("/:disassembly-template-id", controller.DisassemblyTemplateController.Get)       //获取项目拆解模板详情
			disassemblyTemplate.POST("", controller.DisassemblyTemplateController.Create)                            //新增项目拆解模板
			disassemblyTemplate.PUT("/:disassembly-template-id", controller.DisassemblyTemplateController.Update)    //修改项目拆解模板
			disassemblyTemplate.DELETE("/:disassembly-template-id", controller.DisassemblyTemplateController.Delete) //删除项目拆解模板
			disassemblyTemplate.POST("/list", controller.DisassemblyTemplateController.List)                         //获取项目拆解模板列表
		}
		operationRecord := api.Group("/operation-record")
		{
			operationRecord.GET("/:operation-record-id", controller.OperationRecordController.Get)       //获取操作记录详情
			operationRecord.POST("", controller.OperationRecordController.Create)                        //新增操作记录
			operationRecord.PUT("/:operation-record-id", controller.OperationRecordController.Update)    //修改操作记录
			operationRecord.DELETE("/:operation-record-id", controller.OperationRecordController.Delete) //删除操作记录
			operationRecord.POST("/list", controller.OperationRecordController.List)                     //获取操作详情列表
		}
		errorLog := api.Group("/error-log")
		{
			errorLog.GET("/:error-log-id", controller.ErrorLogController.Get)       //获取错误日志详情
			errorLog.POST("", controller.ErrorLogController.Create)                 //新增错误日志
			errorLog.PUT("/:error-log-id", controller.ErrorLogController.Update)    //修改错误日志
			errorLog.DELETE("/:error-log-id", controller.ErrorLogController.Delete) //删除错误日志

		}
		//数据字典的类型
		dictionaryType := api.Group("/dictionary-type")
		{
			dictionaryType.GET("/:dictionary-type-id", controller.DictionaryTypeController.Get)                                  //获取字典类型
			dictionaryType.POST("", middleware.OperationLog(), controller.DictionaryTypeController.Create)                       //新增字典类型
			dictionaryType.POST("/batch", middleware.OperationLog(), controller.DictionaryTypeController.CreateInBatches)        //批量新增字典类型
			dictionaryType.PATCH("/:dictionary-type-id", middleware.OperationLog(), controller.DictionaryTypeController.Update)  //修改字典类型
			dictionaryType.DELETE("/:dictionary-type-id", middleware.OperationLog(), controller.DictionaryTypeController.Delete) //删除字典类型
			dictionaryType.POST("/array", controller.DictionaryTypeController.GetArray)                                          //获取字典类型的数组
			dictionaryType.POST("/list", controller.DictionaryTypeController.GetList)                                            //获取字典类型的列表
		}
		//数据字典的详情项
		dictionaryItem := api.Group("/dictionary-item")
		{
			dictionaryItem.GET("/:dictionary-item-id", controller.DictionaryItemController.Get)       //获取字典项的值
			dictionaryItem.POST("", controller.DictionaryItemController.Create)                       //新增字典项的值
			dictionaryItem.POST("/batch", controller.DictionaryItemController.CreateInBatches)        //批量新增字典项的值
			dictionaryItem.PATCH("/:dictionary-item-id", controller.DictionaryItemController.Update)  //修改字典项的值
			dictionaryItem.DELETE("/:dictionary-item-id", controller.DictionaryItemController.Delete) //删除字典项的值
			dictionaryItem.POST("/array", controller.DictionaryItemController.GetArray)               //获取字典项的数组
			dictionaryItem.POST("/list", controller.DictionaryItemController.GetList)                 //获取字典项的列表
		}
	}

	//引擎配置完成后，返回
	return engine
}
