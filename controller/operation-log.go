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

type operationLog struct{}

func (*operationLog) Get(c *gin.Context) {
	operationLogID, err := strconv.Atoi(c.Param("operation-log-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.OperationRecord.Get(operationLogID)
	c.JSON(http.StatusOK, res)
	return
}

func (*operationLog) Create(c *gin.Context) {
	var param dto.OperationLogCreateOrUpdate
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	//处理lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.OperationRecord.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (*operationLog) Update(c *gin.Context) {
	var param dto.OperationLogCreateOrUpdate
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("operation-log-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		global.SugaredLogger.Errorln(err)
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.OperationRecord.Update(&param)
	c.JSON(200, res)
}

func (*operationLog) Delete(c *gin.Context) {
	operationRecordID, err := strconv.Atoi(c.Param("operation-log-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.OperationRecord.Delete(operationRecordID)
	c.JSON(http.StatusOK, res)
}

func (*operationLog) List(c *gin.Context) {
	var param dto.OperationRecordList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := service.OperationRecord.List(param)
	c.JSON(http.StatusOK, res)
}
