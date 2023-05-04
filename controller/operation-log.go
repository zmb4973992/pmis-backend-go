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

type operationLog struct{}

func (o *operationLog) Get(c *gin.Context) {
	var param service.OperationLogGet
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("operation-log-id"), 10, 64)
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

func (o *operationLog) Delete(c *gin.Context) {
	var param service.OperationLogDelete
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("operation-log-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := param.Delete()
	c.JSON(http.StatusOK, res)
	return
}

func (o *operationLog) GetList(c *gin.Context) {
	var param service.OperationLogGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}
