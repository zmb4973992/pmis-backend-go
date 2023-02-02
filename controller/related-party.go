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

type relatedParty struct{}

func (*relatedParty) Get(c *gin.Context) {
	relatedPartyID, err := strconv.Atoi(c.Param("related-party-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.RelatedParty.Get(relatedPartyID)
	c.JSON(http.StatusOK, res)
	return
}

func (*relatedParty) Create(c *gin.Context) {
	var param dto.RelatedPartyCreateOrUpdate
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.RelatedParty.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (*relatedParty) Update(c *gin.Context) {
	var param dto.RelatedPartyCreateOrUpdate
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("related-party-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.RelatedParty.Update(&param)
	c.JSON(200, res)
}

func (*relatedParty) Delete(c *gin.Context) {
	relatedPartyID, err := strconv.Atoi(c.Param("related-party-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}
	res := service.RelatedParty.Delete(relatedPartyID)
	c.JSON(http.StatusOK, res)
}

func (*relatedParty) List(c *gin.Context) {
	var param dto.RelatedPartyList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成userService,然后调用它的方法
	res := service.RelatedParty.List(param)
	c.JSON(http.StatusOK, res)
}
