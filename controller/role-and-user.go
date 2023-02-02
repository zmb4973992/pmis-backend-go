package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type roleAndUser struct{}

func (roleAndUser) ListByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.ListByRoleID(roleID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) CreateByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.UserIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.CreateByRoleID(roleID, param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) UpdateByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.UserIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	//处理lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.UpdateByRoleID(roleID, param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) DeleteByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.DeleteByRoleID(roleID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) ListByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.ListByUserID(userID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) CreateByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.RoleIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.CreateByUserID(userID, param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) UpdateByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.RoleIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	//处理lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.UpdateByUserID(userID, param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) DeleteByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.DeleteByUserID(userID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUser) ListByTokenInHeader(c *gin.Context) {
	//通过中间件，设定header必须带有token才能访问
	//header里有token后，中间件会自动在context里添加user_id属性，详见自定义的中间件
	tempUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorAccessTokenInvalid))
		return
	}
	userID := tempUserID.(int)
	res := service.RoleAndUser.ListByUserID(userID)
	c.JSON(http.StatusOK, res)
}
