package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

func Init() {
	//检查上传文件的文件夹是否存在
	exists := util.PathExistsOrNot(global.Config.StoragePath)
	//如果不存在就创建
	if !exists {
		err := os.MkdirAll(global.Config.StoragePath, os.ModePerm)
		if err != nil {
			global.SugaredLogger.Panicln(err)
		}
	}
}

type file struct{}

func (f *file) Get(c *gin.Context) {
	var param service.FileGet
	var err error
	param.ID, err = strconv.ParseInt(c.Param("file-id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	filePath, fileName, existed := param.Get()
	if !existed {
		c.JSON(http.StatusOK, response.Failure(util.ErrorFileNotFound))
	}

	c.FileAttachment(filePath, fileName)
	return
}

func (f *file) Create(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorFailToUploadFiles))
		return
	}

	var param service.FileCreate
	param.FileHeader = fileHeader

	userID, exists := util.GetUserID(c)
	if exists {
		param.UserID = userID
	}

	id, url, err1 := param.Create()
	if err1 != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorFailToUploadFiles))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(
		gin.H{
			"id":  id,
			"url": url,
		},
	))
	return

}

func (f *file) Delete(c *gin.Context) {
	var param service.FileDelete
	var err error
	param.ID, err = strconv.ParseInt(c.Param("file-id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	res := param.Delete()
	c.JSON(http.StatusOK, res)
	return
}
