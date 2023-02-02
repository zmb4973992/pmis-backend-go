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

type department struct{}

func (*department) Get(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("department-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.Department.Get(departmentID)
	c.JSON(http.StatusOK, res)
}

func (*department) Create(c *gin.Context) {
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

	res := service.Department.Create(param)
	c.JSON(http.StatusOK, res)
}

func (*department) Update(c *gin.Context) {
	var param dto.DepartmentUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	param.ID, err = strconv.Atoi(c.Param("department-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := c.Get("user_id")
	if exists {
		param.LastModifier = userID.(int)
	}

	res := service.Department.Update(param)
	c.JSON(http.StatusOK, res)
}

func (*department) Delete(c *gin.Context) {
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
	res := service.Department.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (*department) GetArray(c *gin.Context) {
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
	res := service.Department.GetArray(param)
	c.JSON(http.StatusOK, res)
}

func (*department) GetList(c *gin.Context) {
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
	res := service.Department.List(param)
	c.JSON(http.StatusOK, res)
}
