package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"learn-go/dto"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/service"
	"learn-go/util"
	"net/http"
	"strconv"
)

/* controller层负责接收参数、校验参数
然后把id或dto传给service层进行业务处理
最后拿到service层返回的结果进行展现
*/

type relatedPartyController struct{}

func (relatedPartyController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.RelatedPartyService.Get(id)
	c.JSON(http.StatusOK, res)
	return
}

func (relatedPartyController) Update(c *gin.Context) {
	var paramIn model.RelatedParty
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&paramIn)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	paramIn.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.RelatedPartyService.Update(&paramIn)
	c.JSON(200, res)
}

func (relatedPartyController) Create(c *gin.Context) {
	var paramIn model.RelatedParty
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&paramIn)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.RelatedPartyService.Create(&paramIn)
	c.JSON(http.StatusOK, res)
	return
}

func (relatedPartyController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.RelatedPartyService.Delete(id)
	c.JSON(http.StatusOK, res)
}

func (relatedPartyController) List(c *gin.Context) {
	var param dto.RelatedPartyListDTO
	err := c.ShouldBindJSON(&param)
	if err != nil && errors.Is(err, io.EOF) == false {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}
	//生成userService,然后调用它的方法
	res := service.RelatedPartyService.List(param)
	c.JSON(http.StatusOK, res)
}
