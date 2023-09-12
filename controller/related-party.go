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

type relatedParty struct{}

func (r *relatedParty) Get(c *gin.Context) {
	var param service.RelatedPartyGet
	var err error
	param.Id, err = strconv.ParseInt(c.Param("related-party-id"), 10, 64)
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

func (r *relatedParty) Create(c *gin.Context) {
	var param service.RelatedPartyCreate
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
		param.LastModifier = userId
	}

	errCode := param.Create()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}

func (r *relatedParty) Update(c *gin.Context) {
	var param service.RelatedPartyUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidJSONParameters),
		)
		return
	}

	param.Id, err = strconv.ParseInt(c.Param("related-party-id"), 10, 64)
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

func (r *relatedParty) Delete(c *gin.Context) {
	var param service.RelatedPartyDelete
	var err error
	param.Id, err = strconv.ParseInt(c.Param("related-party-id"), 10, 64)
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

func (r *relatedParty) GetList(c *gin.Context) {
	var param service.RelatedPartyGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest,
			response.GenerateList(nil, util.ErrorInvalidJSONParameters, nil))
		return
	}

	outputs, errCode, paging := param.GetList()
	c.JSON(
		http.StatusOK,
		response.GenerateList(outputs, errCode, paging),
	)
	return
}
