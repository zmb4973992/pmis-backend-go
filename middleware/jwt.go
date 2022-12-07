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

//SetUserInfo 在请求中设置userID、roleNames、departmentIDs、departmentNames
func SetUserInfo(c *gin.Context, userID int) {
	c.Set("user_id", userID)
	var user model.User
	//预加载关联的全部子表信息
	global.DB.Where("id = ?", userID).Preload(clause.Associations).First(&user)

	//设置角色名称数组
	var roleNames []string
	for _, role := range user.Roles {
		var roleInfo model.Role
		global.DB.Where("id = ?", role.RoleID).First(&roleInfo)
		roleNames = append(roleNames, roleInfo.Name)
	}
	c.Set("role_names", roleNames)

	//设置当前用户最高级别的角色
	//var topRole string
	if util.IsInSlice("管理员", roleNames) {
		//topRole = "管理员"
	} else if util.IsInSlice("公司级", roleNames) {
		//topRole = "公司级"
	} else if util.IsInSlice("事业部级", roleNames) {
		//topRole = "事业部级"

		//设置部门id数组
		var tempBusinessDivisionIDs []int
		global.DB.Model(&model.DepartmentAndUser{}).Where("user_id = ?", userID).
			Select("department_id").Find(&tempBusinessDivisionIDs)

		//设置事业部id数组
		var businessDivisionIDs []int
		for _, businessDivisionID := range tempBusinessDivisionIDs {
			var count int64
			global.DB.Model(&model.Department{}).Where("id = ?", businessDivisionID).
				Where("level = ?", "事业部").Count(&count)
			if count > 0 {
				businessDivisionIDs = append(businessDivisionIDs, businessDivisionID)
			}
		}
		c.Set("business_division_ids", businessDivisionIDs)

	} else if util.IsInSlice("部门级", roleNames) {
		//topRole = "部门级"

		//设置部门id数组
		var tempDepartmentIDs []int
		global.DB.Model(&model.DepartmentAndUser{}).Where("user_id = ?", userID).
			Select("department_id").Find(&tempDepartmentIDs)

		//校验level是否为部门
		var departmentIDs []int
		for _, departmentID := range tempDepartmentIDs {
			var count int64
			global.DB.Model(&model.Department{}).Where("id = ?", departmentID).
				Where("level = ?", "部门").Count(&count)
			if count > 0 {
				departmentIDs = append(departmentIDs, departmentID)
			}
		}

		c.Set("department_ids", departmentIDs)
	} else {
		//topRole = "项目级"
	}
	//c.Set("top_role", topRole)
}
