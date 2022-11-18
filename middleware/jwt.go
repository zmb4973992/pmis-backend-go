package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"pmis-backend-go/util/jwt"
	"time"
)

func NeedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("access_token")
		//如果请求头没有携带access_token
		if token == "" {
			c.JSON(http.StatusOK, response.Failure(util.ErrorAccessTokenNotFound))
			c.Abort()
			return
		}
		//开始校验access_token

		res, err := jwt.ParseToken(token)
		//如果存在错误或token已过期
		if err != nil || res.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusOK, response.Failure(util.ErrorAccessTokenInvalid))
			c.Abort()
			return
		}
		//如果access_token校验通过
		SetUserInfo(c, res.UserID)
		c.Next()
		return
	}
}

func SetUserInfo(c *gin.Context, userID int) {
	c.Set("userID", userID)
	var user model.User
	//预加载关联的全部子表信息
	global.DB.Where("id = ?", userID).Preload(clause.Associations).First(&user)

	var roleNames []string
	for _, role := range user.Roles {
		var roleInfo model.Role
		global.DB.Where("id = ?", role.RoleID).First(&roleInfo)
		roleNames = append(roleNames, roleInfo.Name)
	}
	c.Set("roleNames", roleNames)

	//设置所属部门
	var departmentNames []string
	for _, department := range user.Departments {
		var departmentInfo model.Department
		global.DB.Where("id = ?", department.DepartmentID).First(&departmentInfo)
		departmentNames = append(departmentNames, departmentInfo.Name)
	}
	c.Set("departments", departmentNames)
}
