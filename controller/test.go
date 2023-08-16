package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
)

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success())
	return
}
