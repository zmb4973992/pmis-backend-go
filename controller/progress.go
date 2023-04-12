package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type progress struct{}

func (p *progress) Get(c *gin.Context) {
	var param service.ProgressGet
	var err error
	param.ID, err = strconv.Atoi(c.Param("progress-id"))
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

func (p *progress) Create(c *gin.Context) {
	var param service.ProgressCreate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}

	//处理creator、last_modifier字段
	tempUserID, exists := c.Get("user_id")
	if exists {
		userID := tempUserID.(int)
		param.Creator = userID
		param.LastModifier = userID
	}

	res := param.Create()
	c.JSON(http.StatusOK, res)
	return
}

func (p *progress) Update(c *gin.Context) {
	var param service.ProgressUpdate
	err := c.ShouldBindJSON(&param)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("progress-id"))
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理lastModifier字段
	userID, exists := c.Get("user_id")
	if exists {
		param.LastModifier = userID.(int)
	}

	res := param.Update()
	c.JSON(http.StatusOK, res)
	return
}

func (p *progress) Delete(c *gin.Context) {
	var param service.ProgressDelete
	var err error
	param.ID, err = strconv.Atoi(c.Param("progress-id"))
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

func (p *progress) GetList(c *gin.Context) {
	var param service.ProgressGetList
	err := c.ShouldBindJSON(&param)

	//别的类似方法会增加EOF错误判断，这里没有，因为必须传json参数
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	//authInput需要userID
	//userID, exists := c.Get("user_id")
	//if exists {
	//	param.UserID = userID.(int)
	//}

	res := param.GetList()
	c.JSON(http.StatusOK, res)
	return
}
