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
	gin.SetMode(global.Config.AppConfig.AppMode)
	fmt.Println("当前运行模式为：", gin.Mode())

	engine := gin.New()

	//将目录(root变量)下的所有文件设置为静态文件，可以直接访问
	engine.Static("/static-test", global.Config.UploadConfig.StoragePath)
	engine.GET("/test/:disassembly-id", controller.Test)
	engine.GET("/download/:file-name", controller.Download)

	//路由不匹配时的处理
	engine.NoRoute(controller.NoRoute.NoRoute)

	//全局中间件
	engine.Use(middleware.Logger(), gin.Recovery(), middleware.Cors())

	engine.POST("/api/login", middleware.RateLimit(), controller.Login)                         //用户登录
	engine.POST("/api/user", middleware.RateLimit(), controller.User.Create)                    //添加用户
	engine.GET("/api/validate-token/:token", middleware.RateLimit(), controller.Token.Validate) //单独校验token是否有效
	engine.GET("/api/captcha", controller.Captcha.Get)                                          //获取验证码

	//依次加载所有的路由组，以下都需要jwt验证
	api := engine.Group("/api")
	api.Use(middleware.ValidateToken(), middleware.RateLimit())
	{
		user := api.Group("/user")
		{
			user.GET("/:user-id", controller.User.Get)                              //获取用户详情
			user.GET("", controller.User.GetByToken)                                //根据header里的token获取用户详情
			user.PATCH("/:user-id", controller.User.Update)                         //修改用户（目前为全功能，考虑改成：修改用户基本信息）
			user.DELETE("/:user-id", middleware.NeedAuth(), controller.User.Delete) //删除用户
			user.POST("/list", controller.User.List)                                //获取用户列表
		}
		file := api.Group("/file")
		{
			file.POST("/upload/single", controller.FileManagement.UploadSingleFile)      //上传单个文件
			file.POST("/upload/multiple", controller.FileManagement.UploadMultipleFiles) //上传多个文件
			file.DELETE("/delete/:file-uuid", controller.FileManagement.DeleteFile)      //删除单个文件
		}
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
		relatedParty := api.Group("/related-party")
		{
			relatedParty.GET("/:related-party-id", controller.RelatedParty.Get)       //获取相关方详情
			relatedParty.PATCH("/:related-party-id", controller.RelatedParty.Update)  //修改相关方
			relatedParty.POST("", controller.RelatedParty.Create)                     //新增相关方
			relatedParty.DELETE("/:related-party-id", controller.RelatedParty.Delete) //删除相关方
			relatedParty.POST("/list", controller.RelatedParty.GetList)               //获取相关方列表
		}
		department := api.Group("/department")
		{
			department.GET("/:department-id", controller.Department.Get)       //获取部门详情
			department.POST("", controller.Department.Create)                  //新增部门
			department.PATCH("/:department-id", controller.Department.Update)  //修改部门
			department.DELETE("/:department-id", controller.Department.Delete) //删除部门
			department.POST("/array", controller.Department.GetArray)          //获取部门数组
			department.POST("/list", controller.Department.GetList)            //获取部门列表
		}
		project := api.Group("/project")
		{
			project.GET("/:project-id", controller.Project.Get)       //获取项目详情
			project.POST("", controller.Project.Create)               //新增项目
			project.PATCH("/:project-id", controller.Project.Update)  //修改项目
			project.DELETE("/:project-id", controller.Project.Delete) //删除项目
			project.POST("/list", controller.Project.GetList)         //获取项目列表
		}
		disassembly := api.Group("/disassembly")
		{
			disassembly.GET("/:disassembly-id", controller.Disassembly.Get)                                 //获取项目拆解详情
			disassembly.POST("/tree", controller.Disassembly.Tree)                                          //获取项目拆解的节点树
			disassembly.POST("", controller.Disassembly.Create)                                             //新增项目拆解
			disassembly.POST("/batch", controller.Disassembly.CreateInBatches)                              //批量新增项目拆解
			disassembly.PATCH("/:disassembly-id", controller.Disassembly.Update)                            //修改项目拆解
			disassembly.DELETE("/:disassembly-id", controller.Disassembly.Delete)                           //删除项目拆解
			disassembly.DELETE("/with-subitems/:disassembly-id", controller.Disassembly.DeleteWithSubitems) //删除项目拆解（子项一并删除）
			disassembly.POST("/list", controller.Disassembly.GetList)                                       //获取项目拆解列表
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
		//数据字典的类型
		dictionaryType := api.Group("/dictionary-type")
		{
			dictionaryType.GET("/:dictionary-type-id", controller.DictionaryType.Get)                                  //获取字典类型
			dictionaryType.POST("", middleware.OperationLog(), controller.DictionaryType.Create)                       //新增字典类型
			dictionaryType.POST("/batch", middleware.OperationLog(), controller.DictionaryType.CreateInBatches)        //批量新增字典类型
			dictionaryType.PATCH("/:dictionary-type-id", middleware.OperationLog(), controller.DictionaryType.Update)  //修改字典类型
			dictionaryType.DELETE("/:dictionary-type-id", middleware.OperationLog(), controller.DictionaryType.Delete) //删除字典类型
			dictionaryType.POST("/array", controller.DictionaryType.GetArray)                                          //获取字典类型的数组
			dictionaryType.POST("/list", controller.DictionaryType.GetList)                                            //获取字典类型的列表
		}
		//数据字典的详情项
		dictionaryItem := api.Group("/dictionary-item")
		{
			dictionaryItem.GET("/:dictionary-item-id", controller.DictionaryItem.Get)       //获取字典项的值
			dictionaryItem.POST("", controller.DictionaryItem.Create)                       //新增字典项的值
			dictionaryItem.POST("/batch", controller.DictionaryItem.CreateInBatches)        //批量新增字典项的值
			dictionaryItem.PATCH("/:dictionary-item-id", controller.DictionaryItem.Update)  //修改字典项的值
			dictionaryItem.DELETE("/:dictionary-item-id", controller.DictionaryItem.Delete) //删除字典项的值
			dictionaryItem.POST("/array", controller.DictionaryItem.GetArray)               //获取字典项的数组
			dictionaryItem.POST("/list", controller.DictionaryItem.GetList)                 //获取字典项的列表
		}
	}

	//引擎配置完成后，返回
	return engine
}
