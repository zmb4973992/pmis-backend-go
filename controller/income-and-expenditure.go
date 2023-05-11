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

type incomeAndExpenditure struct{}

func (i *incomeAndExpenditure) Get(c *gin.Context) {
	var param service.IncomeAndExpenditureGet
	var err error
	param.ID, err = strconv.ParseInt(c.Param("income-and-expenditure-id"), 10, 64)
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

func (i *incomeAndExpenditure) Create(c *gin.Context) {
	var param service.IncomeAndExpenditureCreate
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
		param.Creator = userID
		param.LastModifier = userID
	}

	res := param.Create()
	c.JSON(http.StatusOK, res)
	return
}

func (i *incomeAndExpenditure) Update(c *gin.Context) {
	var param service.IncomeAndExpenditureUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.ParseInt(c.Param("income-and-expenditure-id"), 10, 64)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理last_modifier字段
	userID, exists := util.GetUserID(c)
	if exists {
		param.LastModifier = userID
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (i *incomeAndExpenditure) Delete(c *gin.Context) {
	var param service.IncomeAndExpenditureDelete
	var err error
	param.ID, err = strconv.ParseInt(c.Param("income-and-expenditure-id"), 10, 64)
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

func (i *incomeAndExpenditure) GetList(c *gin.Context) {
	var param service.IncomeAndExpenditureGetList
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
	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}
