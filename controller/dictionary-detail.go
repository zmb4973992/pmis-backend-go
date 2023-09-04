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
	var param service.DictionaryDetailGet
	var err error
	param.ID, err = strconv.ParseInt(c.Param("dictionary-detail-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest, response.Common{
			Data:    nil,
			Code:    util.ErrorInvalidURIParameters,
			Message: util.GetErrorDescription(util.ErrorInvalidURIParameters),
		})
		return
	}
	output, errCode := param.Get()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(output, errCode),
	)
	return
}

func (d *dictionaryDetail) Create(c *gin.Context) {
	var param service.DictionaryDetailCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}

	//处理creator、last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	errCode := param.Create()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (d *dictionaryDetail) Update(c *gin.Context) {
	var param service.DictionaryDetailUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.ParseInt(c.Param("dictionary-detail-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	errCode := param.Update()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (d *dictionaryDetail) Delete(c *gin.Context) {
	var param service.DictionaryDetailDelete
	var err error
	param.ID, err = strconv.ParseInt(c.Param("dictionary-detail-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	errCode := param.Delete()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (d *dictionaryDetail) GetList(c *gin.Context) {
	var param service.DictionaryDetailGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateList(nil, util.ErrorInvalidJSONParameters, nil),
		)
		return
	}

	outputs, errCode, paging := param.GetList()
	c.JSON(
		http.StatusOK,
		response.GenerateList(outputs, errCode, paging),
	)
	return
}
