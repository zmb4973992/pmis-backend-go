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
}

func (*user) Create(c *gin.Context) {
	var param dto.UserCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Creator = userID
		param.LastModifier = userID
	}

	res := service.User.Create(param)
	c.JSON(http.StatusOK, res)
}

func (*user) Update(c *gin.Context) {
	var param dto.UserUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	param.ID, err = strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := c.Get("user_id")
	if exists {
		param.LastModifier = userID.(int)
	}

	res := service.User.Update(param)
	c.JSON(http.StatusOK, res)
}

func (*user) Delete(c *gin.Context) {
	var param dto.UserDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("user-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理deleter字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Deleter = userID
	}
	res := service.User.Delete(param)
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

	res := service.User.List(param)
	c.JSON(http.StatusOK, res)
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
}
