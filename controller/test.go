package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"strings"
)

func Test(c *gin.Context) {
	res1 := util.GetDataRangeIDs(1)
	c.JSON(http.StatusOK, gin.H{
		"data": res1,
	})
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
