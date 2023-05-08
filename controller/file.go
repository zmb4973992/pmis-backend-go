package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service/upload"
	"pmis-backend-go/util"
	"strconv"
)

type fileManagement struct{}

func (f *fileManagement) UploadSingleFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorFailToUploadFiles))
		return
	}

	oss := upload.NewOss()
	fileName, err := oss.UploadSingleFile(fileHeader)
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorFailToUploadFiles))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(gin.H{
		"file_name": fileName,
	}))
	return
}

func (f *fileManagement) UploadMultipleFiles(c *gin.Context) {
	multiPartForm, err := c.MultipartForm()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorFailToUploadFiles))
		return
	}

	fileHeaders := multiPartForm.File["files"]

	oss := upload.NewOss()
	fileNames, err := oss.UploadMultipleFiles(fileHeaders)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Failure(util.ErrorFailToUploadFiles))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(gin.H{
		"file_names": fileNames,
	}))
	return

}

func (f *fileManagement) DeleteFile(c *gin.Context) {
	fileSnowID, _ := strconv.ParseInt(c.Param("file-snow-id"), 10, 64)

	oss := upload.NewOss()
	err := oss.Delete(fileSnowID)
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorFailToDeleteFiles))
		return
	}

	c.JSON(http.StatusOK, response.Success())
	return

}
