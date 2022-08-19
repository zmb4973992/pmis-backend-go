package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"learn-go/dto"
	"learn-go/serializer/response"
	"learn-go/service"
	"learn-go/util"
	"net/http"
	"strconv"
)

type disassemblyTemplateController struct{}

func (disassemblyTemplateController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DisassemblyTemplateService.Get(id)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyTemplateController) Create(c *gin.Context) {
	var param dto.DisassemblyTemplateCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.DisassemblyTemplateService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyTemplateController) Update(c *gin.Context) {
	var param dto.DisassemblyTemplateCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DisassemblyTemplateService.Update(&param)
	c.JSON(200, res)
}

func (disassemblyTemplateController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DisassemblyTemplateService.Delete(id)
	c.JSON(http.StatusOK, res)
}

func (disassemblyTemplateController) List(c *gin.Context) {
	var param dto.DisassemblyTemplateListDTO
	err := c.ShouldBindJSON(&param)

	if err != nil && errors.Is(err, io.EOF) == false {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := service.DisassemblyTemplateService.List(param)
	c.JSON(http.StatusOK, res)
}
