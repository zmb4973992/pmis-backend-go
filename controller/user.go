package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type user struct{}

func (*user) Get(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.User.Get(userID)
	c.JSON(http.StatusOK, res)
	return
}

func (*user) Create(c *gin.Context) {
	//先声明空的dto，再把context里的数据绑到dto上
	var param dto.UserCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
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

	res := service.User.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

// Update controller的功能：解析uri参数、json参数，拦截非法参数，然后传给service层处理
func (*user) Update(c *gin.Context) {
	//这里只更新传过来的参数，所以采用map形式
	var param dto.UserUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("user-id"))
	//如果解析失败，例如URI的参数不是数字
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	//参数解析完毕，交给service层处理
	res := service.User.Update(&param)
	c.JSON(200, res)
}

func (*user) Delete(c *gin.Context) {
	//把uri上的id参数传递给结构体形式的入参
	userID, err := strconv.Atoi(c.Param("user-id"))
	//如果解析失败，例如URI的参数不是数字
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.User.Delete(userID)
	c.JSON(http.StatusOK, res)
}

func (*user) List(c *gin.Context) {
	var param dto.UserList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}

	//生成userService,然后调用它的方法
	res := service.User.List(param)
	c.JSON(http.StatusOK, res)
	return
}

func (*user) GetByToken(c *gin.Context) {
	//通过中间件，设定header必须带有token才能访问
	//header里有token后，中间件会自动在context里添加user_id属性，详见自定义的中间件
	tempUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorAccessTokenInvalid))
		return
	}
	userID := tempUserID.(int)
	res := service.User.Get(userID)
	c.JSON(http.StatusOK, res)
	return
}
