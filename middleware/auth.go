package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// NeedAuth 如果需要根据角色进行鉴权（casbin进行操作），则使用该中间件
// 使用了这个中间件后，相关请求就会先走casbin的规则
func NeedAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tempUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusOK, response.Fail(util.ErrorUserIDDoesNotExist))
			c.Abort()
			return
		}

		userID := tempUserID.(int)
		var roleIDs []int
		global.DB.Model(&model.RoleAndUser{}).Where("user_id = ?", userID).
			Select("role_id").Find(&roleIDs)
		var roleNames []string
		global.DB.Model(&model.Role{}).Where("id in ?", roleIDs).
			Select("name").Find(&roleNames)
		if len(roleNames) == 0 {
			c.JSON(http.StatusOK, response.Fail(util.ErrorRoleInfoNotFound))
			c.Abort()
			return
		}

		subjects := roleNames        //获取用户角色,casbin规则的主体参数
		object := c.Request.URL.Path //获取请求路径，casbin规则的客体参数
		act := c.Request.Method      //获取请求方法，casbin规则的动作参数
		e := util.NewEnforcer()
		//对角色数组进行遍历
		for _, subject := range subjects {
			//如果角色符合casbin的规则
			ok, _ := e.Enforce(subject, object, act)
			if ok {
				c.Next()
				return
			}
		}
		//循环结束，没有满足条件的角色，则中断请求
		c.JSON(http.StatusOK, response.Fail(util.ErrorRolePermissionDenied))
		c.Abort()
	}

}
