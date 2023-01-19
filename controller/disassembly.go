package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type disassemblyController struct{}

func (disassemblyController) Get(c *gin.Context) {
	disassemblyID, err := strconv.Atoi(c.Param("disassembly-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.DisassemblyService.Get(disassemblyID)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) Tree(c *gin.Context) {
	var param dto.DisassemblyTreeDTO
	err := c.ShouldBindJSON(&param)
	//这里json参数必填，否则无法知道要找哪条记录。因此不能忽略掉EOF错误
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.DisassemblyService.Tree(param)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) Create(c *gin.Context) {
	var param dto.DisassemblyCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.DisassemblyService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) CreateInBatches(c *gin.Context) {
	var param []dto.DisassemblyCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		for i := range param {
			param[i].Creator = &userID
			param[i].LastModifier = &userID
		}
	}

	res := service.DisassemblyService.CreateInBatches(param)
	c.JSON(http.StatusOK, res)
	return
}

func (disassemblyController) Update(c *gin.Context) {
	var param dto.DisassemblyCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("disassembly-id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.DisassemblyService.Update(&param)
	c.JSON(http.StatusOK, res)
}

func (disassemblyController) Delete(c *gin.Context) {
	disassemblyID, err := strconv.Atoi(c.Param("disassembly-id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DisassemblyService.Delete(disassemblyID)
	c.JSON(http.StatusOK, res)
}

func (disassemblyController) DeleteWithSubitems(c *gin.Context) {
	disassemblyID, err := strconv.Atoi(c.Param("disassembly-id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DisassemblyService.DeleteWithSubitems(disassemblyID)
	c.JSON(http.StatusOK, res)
}

func (disassemblyController) List(c *gin.Context) {
	var param dto.DisassemblyListDTO
	err := c.ShouldBindQuery(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := service.DisassemblyService.List(param)
	c.JSON(http.StatusOK, res)
}
