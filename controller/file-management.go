package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service/upload"
	"pmis-backend-go/util"
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
	accessPath, fileName, err := oss.UploadSingleFile(fileHeader)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(util.ErrorFailToUploadFiles))
		return
	}

	c.JSON(http.StatusOK, response.SucceedWithData(gin.H{
		"access_path": accessPath,
		"file_name":   fileName,
	}))
	return
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
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Fail(util.ErrorFailToUploadFiles))
		return
	}

	c.JSON(http.StatusOK, response.SucceedWithData(gin.H{
		"storage_path": storagePath,
		"file_names":   fileNames,
	}))
	return

}

func (f *fileManagement) DeleteFile(c *gin.Context) {

	uuid := c.Param("file-uuid")

	//处理deleter字段
	//tempUserID, exists := c.Get("user_id")
	//if exists {
	//	userID := tempUserID.(int)
	//	param.Deleter = userID
	//}
	//res := param.Delete()
	//
	//if err != nil {
	//	c.JSON(http.StatusOK, response.Fail(util.ErrorFailToDeleteFiles))
	//	return
	//}

	oss := upload.NewOss()
	err := oss.Delete(uuid)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(util.ErrorFailToDeleteFiles))
		return
	}

	c.JSON(http.StatusOK, response.Succeed())
	return

}
