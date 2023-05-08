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

type dictionaryDetail struct{}

func (d *dictionaryDetail) Get(c *gin.Context) {
	param := service.OrganizationGet{}
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("dictionary-detail-snow-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest, response.Common{
			Data:    nil,
			Code:    util.ErrorInvalidURIParameters,
			Message: util.GetMessage(util.ErrorInvalidURIParameters),
		})
		return
	}
	res := param.Get()
	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryDetail) Create(c *gin.Context) {
	var param service.DictionaryDetailCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	userSnowID, exists := util.GetUserSnowID(c)
	if exists {
		param.Creator = userSnowID
		param.LastModifier = userSnowID
	}

	res := param.Create()
	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryDetail) CreateInBatches(c *gin.Context) {
	var param service.DictionaryDetailCreateInBatches
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	userSnowID, exists := util.GetUserSnowID(c)
	if exists {
		for i := range param.Data {
			param.Data[i].Creator = userSnowID
			param.Data[i].LastModifier = userSnowID
		}
	}

	res := param.CreateInBatches()
	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryDetail) Update(c *gin.Context) {
	var param service.DictionaryDetailUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.SnowID, err = strconv.ParseInt(c.Param("dictionary-detail-snow-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userSnowID, exists := util.GetUserSnowID(c)
	if exists {
		param.LastModifier = userSnowID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryDetail) Delete(c *gin.Context) {
	var param service.DictionaryDetailDelete
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("dictionary-detail-snow-id"), 10, 64)
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

func (d *dictionaryDetail) GetList(c *gin.Context) {
	var param service.DictionaryDetailGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}