package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type noRoute struct {
}

func (n *noRoute) NoRoute(c *gin.Context) {
	c.JSON(http.StatusBadRequest,
		response.Failure(util.ErrorInvalidRequest))
	return
}
