package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type noRouteController struct {
}

func (noRouteController) NoRoute(c *gin.Context) {
	c.JSON(http.StatusBadRequest, response.Failure(util.ErrorInvalidRequest))
}
