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

type projectDisassemblyController struct{}

func (projectDisassemblyController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.ProjectDisassemblyService.Get(id)
	c.JSON(http.StatusOK, res)
	return
}

func (projectDisassemblyController) Create(c *gin.Context) {
	var param dto.ProjectDisassemblyCreateAndUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.ProjectDisassemblyService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (projectDisassemblyController) Update(c *gin.Context) {
	var param dto.ProjectDisassemblyCreateAndUpdateDTO
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
	res := service.ProjectDisassemblyService.Update(&param)
	c.JSON(200, res)
}

func (projectDisassemblyController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.ProjectDisassemblyService.Delete(id)
	c.JSON(http.StatusOK, res)
}

func (projectDisassemblyController) List(c *gin.Context) {
	var param dto.ProjectDisassemblyListDTO
	err := c.ShouldBindJSON(&param)

	if err != nil && errors.Is(err, io.EOF) == false {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := service.ProjectDisassemblyService.List(param)
	c.JSON(http.StatusOK, res)
}
