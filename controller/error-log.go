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

type errorLog struct{}

func (*errorLog) Get(c *gin.Context) {
	errorLogID, err := strconv.Atoi(c.Param("error-log-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.ErrorLog.Get(errorLogID)
	c.JSON(http.StatusOK, res)
}

func (*errorLog) Create(c *gin.Context) {
	var param dto.ErrorLogCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
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

	res := service.ErrorLog.Create(param)
	c.JSON(http.StatusOK, res)
}

func (*errorLog) Update(c *gin.Context) {
	var param dto.ErrorLogUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	param.ID, err = strconv.Atoi(c.Param("error-log-id"))
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

	res := service.ErrorLog.Update(param)
	c.JSON(http.StatusOK, res)
}

func (*errorLog) Delete(c *gin.Context) {
	var param dto.ErrorLogDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("error-log-id"))
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
	res := service.ErrorLog.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (*errorLog) GetList(c *gin.Context) {
	var param dto.ErrorLogList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}

	//生成Service,然后调用它的方法
	res := service.ErrorLog.GetList(param)
	c.JSON(http.StatusOK, res)
}
