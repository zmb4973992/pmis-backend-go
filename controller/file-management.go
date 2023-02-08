package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"pmis-backend-go/util/upload"
)

type fileManagement struct{}

func (f *fileManagement) UploadSingleFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorFailToUploadFiles))
		return
	}

	oss := upload.NewOss()
	storagePath, fileName, err := oss.UploadSingleFile(fileHeader)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(util.ErrorFailToUploadFiles))
		return
	}

	c.JSON(http.StatusOK, response.SucceedWithData(gin.H{
		"storage_path": storagePath,
		"file_name":    fileName,
	}))
}

func (f *fileManagement) UploadMultipleFiles(c *gin.Context) {
	multiPartForm, err := c.MultipartForm()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK,
			response.Fail(util.ErrorFailToUploadFiles))
		return
	}

	fileHeaders := multiPartForm.File["files"]

	oss := upload.NewOss()
	storagePath, fileNames, err := oss.UploadMultipleFiles(fileHeaders)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(util.ErrorFailToUploadFiles))
		return
	}

	c.JSON(http.StatusOK, response.SucceedWithData(gin.H{
		"storage_path": storagePath,
		"file_names":   fileNames,
	}))
}
