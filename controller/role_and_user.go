package controller

import (
	"github.com/gin-gonic/gin"
	"learn-go/dto"
	"learn-go/serializer/response"
	"learn-go/service"
	"learn-go/util"
	"net/http"
	"strconv"
)

type roleAndUserController struct{}

func (roleAndUserController) ListByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUserService.ListByRoleID(roleID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) CreateByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.UserIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.RoleAndUserService.CreateByRoleID(roleID, param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) UpdateByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.UserIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.RoleAndUserService.UpdateByRoleID(roleID, param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) DeleteByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUserService.DeleteByRoleID(roleID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) ListByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUserService.ListByUserID(userID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) CreateByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var data dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&data)
	if err != nil || len(data.RoleIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.UserIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.RoleAndUserService.CreateByUserID(userID, data)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) UpdateByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var param dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&param)
	if err != nil || len(param.RoleIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.RoleAndUserService.UpdateByUserID(userID, param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) DeleteByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := service.RoleAndUserService.DeleteByUserID(userID)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) ListByTokenInHeader(c *gin.Context) {
	//通过中间件，设定header必须带有token才能访问
	//header里有token后，中间件会自动在context里添加user_id属性，详见自定义的中间件
	tempUserID, ok := c.Get("user_id")
	if ok == false {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorAccessTokenInvalid))
		return
	}
	userID := tempUserID.(int)
	res := service.RoleAndUserService.ListByUserID(userID)
	c.JSON(http.StatusOK, res)
}
