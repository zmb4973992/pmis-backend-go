package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type noMethodController struct {
}

func (noMethodController) NoMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed,
		response.Fail(util.ErrorMethodNotAllowed))
}
