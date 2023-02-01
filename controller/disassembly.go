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

type disassemblyController struct{}

func (disassemblyController) Get(c *gin.Context) {
	disassemblyID, err := strconv.Atoi(c.Param("disassembly-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.DisassemblyService.Get(disassemblyID)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) Tree(c *gin.Context) {
	var param dto.DisassemblyTree
	err := c.ShouldBindJSON(&param)
	//这里json参数必填，否则无法知道要找哪条记录。因此不能忽略掉EOF错误
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.DisassemblyService.Tree(param)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) Create(c *gin.Context) {
	var param dto.DisassemblyCreate
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

	res := service.DisassemblyService.Create(param)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) CreateInBatches(c *gin.Context) {
	var param []dto.DisassemblyCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		for i := range param {
			param[i].Creator = userID
			param[i].LastModifier = userID
		}
	}

	res := service.DisassemblyService.CreateInBatches(param)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) Update(c *gin.Context) {
	var param dto.DisassemblyUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("disassembly-id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := c.Get("user_id")
	if exists {
		param.LastModifier = userID.(int)
	}

	res := service.DisassemblyService.Update(param)
	c.JSON(http.StatusOK, res)
}

func (disassemblyController) Delete(c *gin.Context) {
	var param dto.DisassemblyDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("disassembly-id"))
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
	res := service.DisassemblyService.Delete(param)
	c.JSON(http.StatusOK, res)
}

func (disassemblyController) DeleteWithSubitems(c *gin.Context) {
	var param dto.DisassemblyDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("disassembly-id"))
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

	res := service.DisassemblyService.DeleteWithSubitems(param)
	c.JSON(http.StatusOK, res)
}

func (disassemblyController) GetList(c *gin.Context) {
	var param dto.DisassemblyList
	err := c.ShouldBindJSON(&param)

	//如果json没有传参，会提示EOF错误，这里允许正常运行(允许不传参的查询)；
	//如果是其他错误，就正常报错
	if err != nil && errors.Is(err, io.EOF) {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailForList(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.DisassemblyService.GetList(param)
	c.JSON(http.StatusOK, res)
}
