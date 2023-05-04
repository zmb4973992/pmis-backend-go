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

type dictionaryType struct {
}

func (d *dictionaryType) Get(c *gin.Context) {
	var param = service.DictionaryTypeGet{}
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("dictionary-type-id"), 10, 64)
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

func (d *dictionaryType) Create(c *gin.Context) {
	var param service.DictionaryTypeCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	userID, exists := util.GetUserSnowID(c)
	if exists {
		param.Creator = userID
		param.LastModifier = userID
	}

	res := param.Create()

	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryType) CreateInBatches(c *gin.Context) {
	var param service.DictionaryTypeCreateInBatches
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	userID, exists := util.GetUserSnowID(c)
	if exists {
		for i := range param.Data {
			param.Data[i].Creator = userID
			param.Data[i].LastModifier = userID
		}
	}

	res := param.CreateInBatches()
	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryType) Update(c *gin.Context) {
	var param service.DictionaryTypeUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	param.SnowID, err = strconv.ParseInt(c.Param("dictionary-type-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserSnowID(c)
	if exists {
		param.LastModifier = userID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryType) Delete(c *gin.Context) {
	var param service.DictionaryTypeDelete
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("dictionary-type-id"), 10, 64)
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

func (d *dictionaryType) GetArray(c *gin.Context) {
	var param service.DictionaryTypeGetArray
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
	res := param.GetArray()
	c.JSON(http.StatusOK, res)
	return
}

func (d *dictionaryType) GetList(c *gin.Context) {
	var param service.DictionaryTypeGetList
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
