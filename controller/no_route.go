package controller

import (
	"github.com/gin-gonic/gin"
	"learn-go/serializer/response"
	"learn-go/util"
	"net/http"
)

type noRouteController struct {
}

func (noRouteController) NoRoute(c *gin.Context) {
	c.JSON(http.StatusBadRequest, response.Failure(util.ErrorInvalidRequest))
}
