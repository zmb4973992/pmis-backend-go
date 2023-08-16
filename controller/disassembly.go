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

type disassembly struct{}

func (d *disassembly) Get(c *gin.Context) {
	var param service.DisassemblyGet
	var err error
	param.ID, err = strconv.ParseInt(c.Param("disassembly-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := param.Get()
	c.JSON(http.StatusOK, res)
	return
}

func (d *disassembly) Tree(c *gin.Context) {
	var param service.DisassemblyTree
	err := c.ShouldBindJSON(&param)
	//这里json参数必填，否则无法知道要找哪条记录。因此不能忽略掉EOF错误
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := param.Tree()
	c.JSON(http.StatusOK, res)
	return
}

func (d *disassembly) Create(c *gin.Context) {
	var param service.DisassemblyCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
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

func (d *disassembly) Update(c *gin.Context) {
	var param service.DisassemblyUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	param.ID, err = strconv.ParseInt(c.Param("disassembly-id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
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

func (d *disassembly) Delete(c *gin.Context) {
	var param service.DisassemblyDelete
	var err error
	param.ID, err = strconv.ParseInt(c.Param("disassembly-id"), 10, 64)
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

func (d *disassembly) GetList(c *gin.Context) {
	var param service.DisassemblyGetList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}
