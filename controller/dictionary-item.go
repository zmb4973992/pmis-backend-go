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

type dictionaryItem struct{}

func (*dictionaryItem) Get(c *gin.Context) {
	dictionaryItemID, err := strconv.Atoi(c.Param("dictionary-item-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest, response.Common{
			Data:    nil,
			Code:    util.ErrorInvalidURIParameters,
			Message: util.GetMessage(util.ErrorInvalidURIParameters),
		})
		return
	}
	res := service.DictionaryItem.Get(dictionaryItemID)
	c.JSON(http.StatusOK, res)
}

func (*dictionaryItem) Create(c *gin.Context) {
	var param dto.DictionaryItemCreate
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

	res := service.DictionaryItem.Create(param)
	c.JSON(http.StatusOK, res)
}

func (*dictionaryItem) CreateInBatches(c *gin.Context) {
	var param []dto.DictionaryItemCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	//处理creator、last_modifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		for i := range param {
			param[i].Creator = userID
			param[i].LastModifier = userID
		}
	}

	res := service.DictionaryItem.CreateInBatches(param)
	c.JSON(http.StatusOK, res)
}

func (*dictionaryItem) Update(c *gin.Context) {
	var param dto.DictionaryItemUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("dictionary-item-id"))
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

	res := service.DictionaryItem.Update(param)
	c.JSON(http.StatusOK, res)
}

func (*dictionaryItem) Delete(c *gin.Context) {
	var param dto.DictionaryItemDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("dictionary-item-id"))
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
	res := service.DictionaryItem.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (*dictionaryItem) GetArray(c *gin.Context) {
	var param dto.DictionaryItemList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.DictionaryItem.GetArray(param)
	c.JSON(http.StatusOK, res)
}

func (*dictionaryItem) GetList(c *gin.Context) {
	var param dto.DictionaryItemList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}

	res := service.DictionaryItem.GetList(param)
	c.JSON(http.StatusOK, res)
}
