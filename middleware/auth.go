package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// NeedAuth 如果需要根据角色进行鉴权（casbin进行操作），则使用该中间件
// 这里通过casbin控制哪些角色可以访问接口、哪些角色不能访问接口
func NeedAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tempSubjects, ok := c.Get("roles") //获取用户角色,casbin规则的主体参数
		subjects := tempSubjects.([]string)
		if !ok || len(subjects) == 0 {
			c.JSON(http.StatusOK, response.Fail(util.ErrorPermissionDenied))
			c.Abort()
			return
		}

		object := c.Request.URL.Path //获取请求路径，casbin规则的客体参数
		act := c.Request.Method      //获取请求方法，casbin规则的动作参数
		e := util.NewEnforcer()
		//对角色列表进行遍历
		for _, subject := range subjects {
			//如果角色符合casbin的规则
			ok, _ := e.Enforce(subject, object, act)
			if ok {
				//放行，跳出循环
				c.Next()
				return
			}
		}
		//循环结束，没有满足条件的角色，则中断请求
		c.JSON(http.StatusOK, response.Fail(util.ErrorPermissionDenied))
		c.Abort()
		return
	}
}
