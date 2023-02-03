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
	res := service.OperationLog.Get(operationLogID)
	c.JSON(http.StatusOK, res)
}

func (*operationLog) Delete(c *gin.Context) {
	var param dto.OperationLogDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("operation-log-id"))
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
	res := service.OperationLog.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (*operationLog) GetList(c *gin.Context) {
	var param dto.OperationLogList
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
	res := service.OperationLog.GetList(param)
	c.JSON(http.StatusOK, res)
}
