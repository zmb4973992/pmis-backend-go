package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"strconv"
)

// RBAC 如果需要根据角色进行鉴权，则使用该中间件
func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := util.GetUserID(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusOK, response.Failure(util.ErrorUserIDDoesNotExist))
			return
		}

		var roleIDs []int
		global.DB.Model(&model.UserAndRole{}).Where("user_id = ?", userID).
			Select("role_id").Find(&roleIDs)

		object := c.Request.URL.Path //获取请求路径，casbin规则的客体参数
		act := c.Request.Method      //获取请求方法，casbin规则的动作参数

		cachedEnforcer, err := util.NewCachedEnforcer()
		if err != nil {
			global.SugaredLogger.Errorln(err)
		}

		var permitted bool
		for _, role := range roleIDs {
			subject := strconv.Itoa(role)
			permitted, _ = cachedEnforcer.Enforce(subject, object, act)
			if permitted {
				break
			}
		}

		if !permitted {
			c.AbortWithStatusJSON(http.StatusOK, response.Failure(util.ErrorUnauthorized))
			return
		}

		c.Next()
		return
	}
}
