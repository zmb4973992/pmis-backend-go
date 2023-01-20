package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type dictionaryTypeController struct{}

func (dictionaryTypeController) Create(c *gin.Context) {
	var param dto.DictionaryTypeCreateOrUpdateDTO
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

	res := service.DictionaryTypeService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (dictionaryTypeController) CreateInBatches(c *gin.Context) {
	var param []dto.DictionaryTypeCreateOrUpdateDTO
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

	res := service.DictionaryTypeService.CreateInBatches(param)
	c.JSON(http.StatusOK, res)
	return
}

func (dictionaryTypeController) Update(c *gin.Context) {
	var param dto.DictionaryTypeCreateOrUpdateDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("dictionary-type-id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.DictionaryTypeService.Update(&param)
	c.JSON(http.StatusOK, res)
}

func (dictionaryTypeController) Delete(c *gin.Context) {
	dictionaryTypeID, err := strconv.Atoi(c.Param("dictionary-type-id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DictionaryTypeService.Delete(dictionaryTypeID)
	c.JSON(http.StatusOK, res)
}

func (dictionaryTypeController) List(c *gin.Context) {
	var param dto.DictionaryTypeListDTO
	err := c.ShouldBindJSON(&param)

	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := service.DictionaryTypeService.List(param)
	c.JSON(http.StatusOK, res)
}
