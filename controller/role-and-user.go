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

func (r *roleAndUser) ListByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.ListByRoleID(roleID)
	c.JSON(http.StatusOK, res)
	return
}

func (r *roleAndUser) CreateByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.UserIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.CreateByRoleID(roleID, param)
	c.JSON(http.StatusOK, res)
	return
}

func (r *roleAndUser) UpdateByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.UserIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.UpdateByRoleID(roleID, param)
	c.JSON(http.StatusOK, res)
	return
}

func (r *roleAndUser) DeleteByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.DeleteByRoleID(roleID)
	c.JSON(http.StatusOK, res)
}

func (r *roleAndUser) ListByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.ListByUserID(userID)
	c.JSON(http.StatusOK, res)
	return
}

func (r *roleAndUser) CreateByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.RoleIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.CreateByUserID(userID, param)
	c.JSON(http.StatusOK, res)
	return
}

func (r *roleAndUser) UpdateByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdate
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.RoleIDs) == 0 {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.LastModifier = &userID
	}

	res := service.RoleAndUser.UpdateByUserID(userID, param)
	c.JSON(http.StatusOK, res)
	return
}

func (r *roleAndUser) DeleteByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUser.DeleteByUserID(userID)
	c.JSON(http.StatusOK, res)
	return
}

func (r *roleAndUser) ListByTokenInHeader(c *gin.Context) {
	//通过中间件，设定header必须带有token才能访问
	//header里有token后，中间件会自动在context里添加user_id属性，详见自定义的中间件
	tempUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorAccessTokenInvalid))
		return
	}
	userID := tempUserID.(int)
	res := service.RoleAndUser.ListByUserID(userID)
	c.JSON(http.StatusOK, res)
	return
}
