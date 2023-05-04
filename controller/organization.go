package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type organization struct{}

func (o *organization) Get(c *gin.Context) {
	param := service.OrganizationGet{}
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("organization-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := param.Get()
	c.JSON(http.StatusOK, res)
	return
}

func (o *organization) Create(c *gin.Context) {
	var param service.OrganizationCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	userID, exists := util.GetUserSnowID(c)
	if exists {
		param.Creator = userID
		param.LastModifier = userID
	}

	res := param.Create()
	c.JSON(http.StatusOK, res)
	return
}

func (o *organization) Update(c *gin.Context) {
	var param service.OrganizationUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	param.SnowID, err = strconv.ParseInt(c.Param("organization-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserSnowID(c)
	if exists {
		param.LastModifier = userID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (o *organization) Delete(c *gin.Context) {
	var param service.OrganizationDelete
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("organization-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := param.Delete()
	c.JSON(http.StatusOK, res)
	return
}

//func (o *organization) GetArray(c *gin.Context) {
//	var param service.OrganizationGetArray
//	err := c.ShouldBindJSON(&param)
//
//	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
//	//如果是其他错误，就正常报错
//	if err != nil && !errors.Is(err, io.EOF) {
//		global.SugaredLogger.Errorln(err)
//		c.JSON(http.StatusBadRequest,
//			response.FailureForList(util.ErrorInvalidJSONParameters))
//		return
//	}
//
//	tempUserID, exists := c.Get("user_id")
//	if exists {
//		userID := tempUserID.(int)
//		param.UserID = userID
//	}
//
//	res := param.GetArray()
//	c.JSON(http.StatusOK, res)
//	return
//}

func (o *organization) GetList(c *gin.Context) {
	var param service.OrganizationGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	//AuthorityInput需要userID
	userID, exists := util.GetUserSnowID(c)
	if exists {
		param.UserSnowID = userID
	}

	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}
