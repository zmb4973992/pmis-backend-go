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
	exists := util.PathExistsOrNot(global.Config.Upload.StoragePath)
	//如果不存在就创建
	if !exists {
		err := os.MkdirAll(global.Config.Upload.StoragePath, os.ModePerm)
		if err != nil {
			global.SugaredLogger.Panicln(err)
		}
	}
}

type file struct{}

func (f *file) Get(c *gin.Context) {
	var param service.FileGet
	var err error
	param.Id, err = strconv.ParseInt(c.Param("file-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	filePath, fileName, existed := param.Get()
	if !existed {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorFileNotFound),
		)
	}

	c.FileAttachment(filePath, fileName)
	return
}

func (f *file) Create(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorFailToUploadFiles),
		)
		return
	}

	var param service.FileCreate
	param.FileHeader = fileHeader

	userId, exists := util.GetUserId(c)
	if exists {
		param.UserId = userId
	}

	id, url, err1 := param.Create()
	if err1 != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorFailToUploadFiles),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response.GenerateCommon(
			gin.H{
				"id":  id,
				"url": url,
			},
			util.Success,
		))
	return
}

func (f *file) Delete(c *gin.Context) {
	var param service.FileDelete
	var err error
	param.Id, err = strconv.ParseInt(c.Param("file-id"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	errCode := param.Delete()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(nil, errCode),
	)
	return
}
