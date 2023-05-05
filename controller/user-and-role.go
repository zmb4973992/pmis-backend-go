package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type userAndRole struct{}

func (u *userAndRole) UpdateByRoleSnowID(c *gin.Context) {
	var param service.RoleAndUserUpdateByRoleSnowID
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.RoleSnowID, err = strconv.ParseInt(c.Param("role-snow-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userSnowID, exists := util.GetUserSnowID(c)
	if exists {
		param.Creator = userSnowID
		param.LastModifier = userSnowID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (u *userAndRole) UpdateByUserSnowID(c *gin.Context) {
	var param service.RoleAndUserUpdateByUserSnowID
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.UserSnowID, err = strconv.ParseInt(c.Param("user-snow-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userSnowID, exists := util.GetUserSnowID(c)
	if exists {
		param.Creator = userSnowID
		param.LastModifier = userSnowID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}
