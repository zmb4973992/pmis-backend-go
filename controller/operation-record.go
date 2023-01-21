package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type operationRecordController struct{}

func (operationRecordController) Get(c *gin.Context) {
	operationRecordID, err := strconv.Atoi(c.Param("operation-record-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.OperationRecordService.Get(operationRecordID)
	c.JSON(http.StatusOK, res)
	return
}

func (operationRecordController) Create(c *gin.Context) {
	var param dto.OperationRecordCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理lastModifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.OperationRecordService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (operationRecordController) Update(c *gin.Context) {
	var param dto.OperationRecordCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("operation-record-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		global.SugaredLogger.Errorln(err)
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.OperationRecordService.Update(&param)
	c.JSON(200, res)
}

func (operationRecordController) Delete(c *gin.Context) {
	operationRecordID, err := strconv.Atoi(c.Param("operation-record-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.OperationRecordService.Delete(operationRecordID)
	c.JSON(http.StatusOK, res)
}

func (operationRecordController) List(c *gin.Context) {
	var param dto.OperationRecordListDTO
	err := c.ShouldBindQuery(&param)

	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成Service,然后调用它的方法
	res := service.OperationRecordService.List(param)
	c.JSON(http.StatusOK, res)
}
