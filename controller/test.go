package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

func Test(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, util.Success),
	)
	return
}
