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

	var data dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&data)
	if err != nil || len(data.UserIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.RoleAndUserService.CreateByRoleID(roleID, data)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) UpdateByRoleID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	var data dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&data)
	if err != nil || len(data.UserIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.RoleAndUserService.UpdateByRoleID(roleID, data)
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

	var data dto.RoleAndUserCreateOrUpdateDTO
	err = c.ShouldBindJSON(&data)
	if err != nil || len(data.RoleIDs) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.RoleAndUserService.UpdateByUserID(userID, data)
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