package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type userController struct{}

func (userController) Get(c *gin.Context) {
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

func (userController) Create(c *gin.Context) {
	//先声明空的dto，再把context里的数据绑到dto上
	var param dto.UserCreateDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
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

	res := service.UserService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

// Update controller的功能：解析uri参数、json参数，拦截非法参数，然后传给service层处理
func (userController) Update(c *gin.Context) {
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

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	//参数解析完毕，交给service层处理
	res := service.UserService.Update(&param)
	c.JSON(200, res)
}

func (userController) Delete(c *gin.Context) {
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

func (userController) List(c *gin.Context) {
	var param dto.UserListDTO
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	//生成userService,然后调用它的方法
	res := service.UserService.List(param)
	c.JSON(http.StatusOK, res)
	return
}

func (userController) GetByToken(c *gin.Context) {
	//通过中间件，设定header必须带有token才能访问
	//header里有token后，中间件会自动在context里添加user_id属性，详见自定义的中间件
	tempUserID, ok := c.Get("user_id")
	if ok == false {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorAccessTokenInvalid))
		return
	}
	userID := tempUserID.(int)
	res := service.UserService.Get(userID)
	c.JSON(http.StatusOK, res)
	return
}
