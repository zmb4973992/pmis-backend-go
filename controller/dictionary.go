package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/service"
)

type dictionaryController struct{}

func (dictionaryController) Get(c *gin.Context) {
	//生成Service,然后调用它的方法
	res := service.DictionaryService.Get()
	c.JSON(http.StatusOK, res)
}
