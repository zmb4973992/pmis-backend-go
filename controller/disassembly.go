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
	param.SnowID, err = strconv.ParseInt(c.Param("disassembly-id"), 10, 64)
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
	userID, exists := util.GetUserSnowID(c)
	if exists {
		param.Creator = userID
		param.LastModifier = userID
	}

	res := param.Create()
	c.JSON(http.StatusOK, res)
	return
}

// deprecated
// 逻辑较复杂，暂时废弃
//func (d *disassembly) CreateInBatches(c *gin.Context) {
//	var param service.DisassemblyCreateInBatches
//	err := c.ShouldBindJSON(&param)
//	if err != nil {
//		global.SugaredLogger.Errorln(err)
//		c.JSON(http.StatusBadRequest,
//			response.Failure(util.ErrorInvalidJSONParameters))
//		return
//	}
//
//	//处理creator、last_modifier字段
//	userID, exists := util.GetUserSnowID(c)
//	if exists {
//		for i := range param.Data {
//			param.Data[i].Creator = userID
//			param.Data[i].LastModifier = userID
//		}
//	}
//
//	res := param.CreateInBatches()
//	c.JSON(http.StatusOK, res)
//	return
//}

func (d *disassembly) Update(c *gin.Context) {
	var param service.DisassemblyUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	param.SnowID, err = strconv.ParseInt(c.Param("disassembly-id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
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

func (d *disassembly) Delete(c *gin.Context) {
	var param service.DisassemblyDelete
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("disassembly-id"), 10, 64)
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

func (d *disassembly) DeleteWithInferiors(c *gin.Context) {
	var param service.DisassemblyDeleteWithInferiors
	var err error
	param.SnowID, err = strconv.ParseInt(c.Param("disassembly-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := param.DeleteWithInferiors()
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
