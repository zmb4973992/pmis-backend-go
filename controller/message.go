package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type message struct{}

func (m *message) Get(c *gin.Context) {
	param := service.MessageGet{}
	var err error
	param.ID, err = strconv.ParseInt(c.Param("message-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := param.Get()
	c.JSON(http.StatusOK, res)
	return
}

func (m *message) Create(c *gin.Context) {
	var param service.MessageCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	println(param.Recipients)

	//处理creator、last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.Create()
	c.JSON(http.StatusOK, res)
	return
}

func (m *message) Update(c *gin.Context) {
	var param service.MessageUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.ParseInt(c.Param("message-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (m *message) Delete(c *gin.Context) {
	var param service.MessageDelete
	var err error
	param.ID, err = strconv.ParseInt(c.Param("message-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.Delete()
	c.JSON(http.StatusOK, res)
	return
}

func (m *message) GetList(c *gin.Context) {
	var param service.MessageGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}