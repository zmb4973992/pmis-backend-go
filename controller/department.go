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

type departmentController struct{}

func (departmentController) Get(c *gin.Context) {
	departmentGet := new(service.DepartmentGet)
	var err error
	departmentGet.DepartmentID, err = strconv.Atoi(c.Param("department-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	//res := service.DepartmentService.Get(departmentID)

	res := departmentGet.Do()
	c.JSON(http.StatusOK, res)

}

func (departmentController) Create(c *gin.Context) {
	var param dto.DepartmentCreate
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

	res := service.DepartmentService.Create(param)
	c.JSON(http.StatusOK, res)
}

func (departmentController) Update(c *gin.Context) {
	var param dto.DepartmentUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("department-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理lastModifier字段
	userID, exists := c.Get("user_id")
	if exists {
		param.LastModifier = userID.(int)
	}

	res := service.DepartmentService.Update(param)
	c.JSON(http.StatusOK, res)
}

func (departmentController) Delete(c *gin.Context) {
	var param dto.DepartmentDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("department-id"))
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
	res := service.DepartmentService.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (departmentController) GetArray(c *gin.Context) {
	var param dto.DepartmentList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}

	//把userID传给dto，给service调用
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.UserID = userID
	}

	//生成Service,然后调用它的方法
	res := service.DepartmentService.GetArray(param)
	c.JSON(http.StatusOK, res)
}

func (departmentController) GetList(c *gin.Context) {
	var param dto.DepartmentList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}

	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.UserID = userID
	}

	//生成Service,然后调用它的方法
	res := service.DepartmentService.List(param)
	c.JSON(http.StatusOK, res)
}
