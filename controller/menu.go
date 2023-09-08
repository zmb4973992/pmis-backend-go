package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type menu struct{}

func (m *menu) Get(c *gin.Context) {
	var param service.MenuGet
	var err error
	param.ID, err = strconv.ParseInt(c.Param("menu-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}
	res := param.Get()
	c.JSON(http.StatusOK, res)
	return
}

func (m *menu) Create(c *gin.Context) {
	var param service.MenuCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}

	//处理creator、last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.Create()
	c.JSON(http.StatusOK, res)
	return
}

func (m *menu) Update(c *gin.Context) {
	var param service.MenuUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.ParseInt(c.Param("menu-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (m *menu) Delete(c *gin.Context) {
	var param service.MenuDelete
	var err error
	param.ID, err = strconv.ParseInt(c.Param("menu-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	res := param.Delete()
	c.JSON(http.StatusOK, res)
	return
}

func (m *menu) GetList(c *gin.Context) {
	var param service.MenuGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateList(nil, util.ErrorInvalidJSONParameters, nil),
		)
		return
	}

	//AuthorityInput需要userID
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}

func (m *menu) GetTree(c *gin.Context) {
	var param service.MenuGetTree
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateList(nil, util.ErrorInvalidJSONParameters, nil),
		)
		return
	}

	//AuthorityInput需要userID
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.GetTree()
	c.JSON(http.StatusOK, res)
	return
}

func (m *menu) UpdateUsers(c *gin.Context) {
	var param service.MenuUpdateApis
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.MenuID, err = strconv.ParseInt(c.Param("menu-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.Creator = userID
		param.LastModifier = userID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}
