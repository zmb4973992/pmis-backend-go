package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

func UploadSingle(c *gin.Context) {
	uniqueFilename, err := util.UploadSingleFile(c, "file")
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorFailToUploadFiles))
		return
	}
	c.JSON(http.StatusOK, response.Common{
		Data: gin.H{
			"unique_filename": uniqueFilename,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	})
	return
}

func UploadMultiple(c *gin.Context) {
	uniqueFilenames, err := util.UploadMultipleFiles(c, "files")
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorFailToUploadFiles))
		return
	}
	c.JSON(http.StatusOK, response.Common{
		Data: gin.H{
			"unique_filenames": uniqueFilenames,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	})
	return
}
