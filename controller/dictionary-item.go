package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type dictionaryItemController struct{}

func (dictionaryItemController) Get(c *gin.Context) {
	dictionaryTypeID, err := strconv.Atoi(c.Param("dictionary-type-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Common{
			Data:    nil,
			Code:    util.ErrorInvalidURIParameters,
			Message: util.GetMessage(util.ErrorInvalidURIParameters),
		})
		return
	}
	//生成Service,然后调用它的方法
	res := service.DictionaryItemService.Get(dictionaryTypeID)
	c.JSON(http.StatusOK, res)
}

func (dictionaryItemController) Create(c *gin.Context) {
	var param dto.DictionaryItemCreateOrUpdateDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.DictionaryItemService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (dictionaryItemController) CreateInBatches(c *gin.Context) {
	var param []dto.DictionaryItemCreateOrUpdateDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		for i := range param {
			param[i].Creator = &userID
			param[i].LastModifier = &userID
		}
	}

	res := service.DictionaryItemService.CreateInBatches(param)
	c.JSON(http.StatusOK, res)
	return
}
