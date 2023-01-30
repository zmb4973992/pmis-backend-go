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

type dictionaryItemController struct{}

func (x dictionaryItemController) Get(c *gin.Context) {
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
	//生成Service,然后调用它的方法
	res := service.DictionaryItemService.Get(dictionaryItemID)
	c.JSON(http.StatusOK, res)
}

func (dictionaryItemController) Create(c *gin.Context) {
	var param dto.DictionaryItemCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//处理creator、last_modifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Creator = userID
		param.LastModifier = userID
	}

	res := service.DictionaryItemService.Create(param)
	c.JSON(http.StatusOK, res)
	return
}

func (dictionaryItemController) CreateInBatches(c *gin.Context) {
	var param []dto.DictionaryItemCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		for i := range param {
			param[i].Creator = userID
			param[i].LastModifier = userID
		}
	}

	res := service.DictionaryItemService.CreateInBatches(param)
	c.JSON(http.StatusOK, res)
	return
}

func (dictionaryItemController) Update(c *gin.Context) {
	var param dto.DictionaryItemUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("dictionary-item-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理lastModifier字段
	userID, exists := c.Get("user_id")
	if exists {
		param.LastModifier = userID.(int)
	}

	res := service.DictionaryItemService.Update(param)
	c.JSON(http.StatusOK, res)
}

func (dictionaryItemController) Delete(c *gin.Context) {
	var param dto.DictionaryItemDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("dictionary-item-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理deleter字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Deleter = userID
	}
	res := service.DictionaryItemService.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (dictionaryItemController) GetArray(c *gin.Context) {
	var param dto.DictionaryItemList
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
	res := service.DictionaryItemService.GetArray(param)
	c.JSON(http.StatusOK, res)
}

func (dictionaryItemController) GetList(c *gin.Context) {
	var param dto.DictionaryItemList
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
	res := service.DictionaryItemService.GetList(param)
	c.JSON(http.StatusOK, res)
}
