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

type user struct{}

func Login(c *gin.Context) {
	var param service.UserLogin
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}

	permitted := param.Verify()
	if !permitted {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorWrongCaptcha),
		)
		return
	}

	output, errCode := param.Login()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(output, errCode),
	)
	return
}

func (u *user) Get(c *gin.Context) {
	var param service.UserGet
	var err error
	param.Id, err = strconv.ParseInt(c.Param("user-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	output, errCode := param.Get()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(output, errCode),
	)
	return
}

func (u *user) Create(c *gin.Context) {
	var param service.UserCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}

	//处理creator、last_modifier字段
	userId, exists := util.GetUserId(c)
	if exists {
		param.UserId = userId
	}

	errCode := param.Create()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (u *user) Update(c *gin.Context) {
	var param service.UserUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}

	param.Id, err = strconv.ParseInt(c.Param("user-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	//处理last_modifier字段
	userId, exists := util.GetUserId(c)
	if exists {
		param.UserId = userId
	}

	errCode := param.Update()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (u *user) Delete(c *gin.Context) {
	var param service.UserDelete
	var err error
	param.Id, err = strconv.ParseInt(c.Param("user-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	errCode := param.Delete()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (u *user) GetList(c *gin.Context) {
	var param service.UserGetList
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

	outputs, errCode, paging := param.GetList()
	c.JSON(
		http.StatusOK,
		response.GenerateList(outputs, errCode, paging),
	)
	return
}

func (u *user) GetByToken(c *gin.Context) {
	var param service.UserGet
	userId, exists := util.GetUserId(c)
	if !exists {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorAccessTokenInvalid),
		)
		return
	}

	param.Id = userId
	output, errCode := param.Get()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(output, errCode),
	)
	return
}

func (u *user) UpdateRoles(c *gin.Context) {
	var param service.UserUpdateRoles
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.UserId, err = strconv.ParseInt(c.Param("user-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	//处理last_modifier字段
	userId, exists := util.GetUserId(c)
	if exists {
		param.LastModifier = userId
	}

	errCode := param.Update()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (u *user) UpdateDataAuthority(c *gin.Context) {
	var param service.UserUpdateDataAuthority
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateList(nil, util.ErrorInvalidJSONParameters, nil),
		)
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.UserId, err = strconv.ParseInt(c.Param("user-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	//处理last_modifier字段
	userId, exists := util.GetUserId(c)
	if exists {
		param.LastModifier = userId
	}

	errCode := param.Update()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}
