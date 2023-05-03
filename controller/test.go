package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strings"
)

func Test(c *gin.Context) {
	var param service.RBACUpdate
	param.ApiIDs = []int{5}

	err := param.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
	return
}

func Download(c *gin.Context) {
	fileName := c.Param("file-name")
	if fileName != "" {
		storagePath := global.Config.UploadConfig.StoragePath
		fileNameForFrontend := strings.Split(fileName, "--")[1]
		_, err := os.Stat(storagePath + fileName)
		if err != nil {
			c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
			return
		}
		c.FileAttachment(storagePath+fileName, fileNameForFrontend)
		return
	}
	c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
	return
}
