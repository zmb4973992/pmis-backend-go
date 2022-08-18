package controller

import (
	"github.com/gin-gonic/gin"
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/serializer/response"
	"learn-go/service"
	"learn-go/util"
	"net/http"
	"strconv"
)

type roleAndUserController struct{}

func (roleAndUserController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.UserService.Get(id)
	c.JSON(http.StatusOK, res)
	return
}

func (roleAndUserController) Create(c *gin.Context) {
	//先声明空的dto，再把context里的数据绑到dto上
	var u dto.UserCreateDTO
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.UserService.Create(&u)
	c.JSON(http.StatusOK, res)
	return
}

// Update controller的功能：解析uri参数、json参数，拦截非法参数，然后传给service层处理
func (roleAndUserController) Update(c *gin.Context) {
	//这里只更新传过来的参数，所以采用map形式
	var param dto.UserUpdateDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("id"))
	//如果解析失败，例如URI的参数不是数字
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	//参数解析完毕，交给service层处理
	res := service.UserService.Update(&param)
	c.JSON(200, res)
}

func (roleAndUserController) Delete(c *gin.Context) {
	//把uri上的id参数传递给结构体形式的入参
	id, err := strconv.Atoi(c.Param("id"))
	//如果解析失败，例如URI的参数不是数字
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.UserService.Delete(id)
	c.JSON(http.StatusOK, res)
}

func (roleAndUserController) UserSlice(c *gin.Context) {
	var param dto.RoleAndUserListDTO
	err := c.ShouldBindJSON(&param)
	if err != nil || param.RoleID == nil {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	res := dao.RoleAndUserDAO.UserSlice(*param.RoleID)
	c.JSON(http.StatusOK, res)
	return
}

func (roleAndUserController) RoleSlice(c *gin.Context) {
	var param dto.RoleAndUserListDTO
	err := c.ShouldBindJSON(&param)
	if err != nil || param.UserID == nil {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	res := dao.RoleAndUserDAO.RoleSlice(*param.UserID)
	c.JSON(http.StatusOK, res)
	return
}

func (roleAndUserController) UpdateUserSlice(c *gin.Context) {

}
