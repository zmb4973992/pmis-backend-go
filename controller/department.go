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

type departmentController struct{}

func (departmentController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Common{
			Data:    nil,
			Code:    util.ErrorInvalidURIParameters,
			Message: util.GetMessage(util.ErrorInvalidURIParameters),
		})
		return
	}
	res := service.DepartmentService.Get(id)
	c.JSON(http.StatusOK, res)
	return
}

func (departmentController) Create(c *gin.Context) {
	var param dto.DepartmentCreateAndUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.DepartmentService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (departmentController) Update(c *gin.Context) {
	var param dto.DepartmentCreateAndUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DepartmentService.Update(&param)
	c.JSON(200, res)
}

func (departmentController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DepartmentService.Delete(id)
	c.JSON(http.StatusOK, res)
}

func (departmentController) List(c *gin.Context) {
	var param dto.DepartmentListDTO
	err := c.ShouldBindJSON(&param)
	//如果json没有传参，会提示EOF错误，这里可以正常运行；如果是其他错误，就正常报错
	if err != nil && errors.Is(err, io.EOF) == false {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := service.DepartmentService.List(param)
	c.JSON(http.StatusOK, res)
}
