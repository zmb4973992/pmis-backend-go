package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"learn-go/dto"
	"learn-go/serializer/response"
	"learn-go/service"
	"learn-go/util"
	"net/http"
	"strconv"
)

type roleAndUserController struct{}

func (roleAndUserController) Create(c *gin.Context) {
	//先声明空的dto，再把context里的数据绑到dto上
	var r dto.RoleAndUserCreateDTO
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.RoleAndUserService.Create(&r)
	c.JSON(http.StatusOK, res)
	return
}

func (roleAndUserController) CreateInBatch(c *gin.Context) {
	var r []dto.RoleAndUserCreateDTO
	err := c.ShouldBindJSON(&r)
	if err != nil || len(r) == 0 {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.RoleAndUserService.CreateInBatch(r)
	c.JSON(http.StatusOK, res)
	return
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
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.RoleAndUserService.UpdateUserIDByRoleID(roleID, data)
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
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.RoleAndUserService.UpdateRoleIDByUserID(userID, data)
	c.JSON(http.StatusOK, res)
}

//	//这里只更新传过来的参数，所以采用map形式
//	var param dto.RoleAndUserCreateDTO
//	err := c.ShouldBindJSON(&param)
//	if err != nil {
//		c.JSON(http.StatusOK,
//			response.Failure(util.ErrorInvalidJSONParameters))
//		return
//	}

func (roleAndUserController) Delete(c *gin.Context) {
	var param dto.RoleAndUserDeleteDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.RoleAndUserService.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) List(c *gin.Context) {
	var param dto.RoleAndUserListDTO
	err := c.ShouldBindJSON(&param)
	if err != nil && errors.Is(err, io.EOF) == false {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成userService,然后调用它的方法
	res := service.RoleAndUserService.List(param)
	c.JSON(http.StatusOK, res)
	return
}

//func (roleAndUserController) UserSlice(c *gin.Context) {
//	var param dto.RoleAndUserListDTO
//	err := c.ShouldBindJSON(&param)
//	if err != nil || param.RoleIDs == nil {
//		c.JSON(http.StatusBadRequest,
//			response.FailureForList(util.ErrorInvalidJSONParameters))
//		return
//	}
//
//	res := dao.RoleAndUserDAO.UserSlice(*param.RoleIDs)
//	c.JSON(http.StatusOK, res)
//	return
//}

//func (roleAndUserController) RoleSlice(c *gin.Context) {
//	var param dto.RoleAndUserListDTO
//	err := c.ShouldBindJSON(&param)
//	if err != nil || param.UserIDs == nil {
//		c.JSON(http.StatusBadRequest,
//			response.FailureForList(util.ErrorInvalidJSONParameters))
//		return
//	}
//
//	res := dao.RoleAndUserDAO.RoleSlice(*param.UserIDs)
//	c.JSON(http.StatusOK, res)
//	return
//}
