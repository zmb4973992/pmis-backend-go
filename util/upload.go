package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"pmis-backend-go/global"
)

func UploadInit() {
	//检查上传文件的文件夹是否存在
	exists := PathExistsOrNot(global.Config.FullPath)
	//如果不存在就创建
	if !exists {
		err := os.MkdirAll(global.Config.FullPath, os.ModePerm)
		if err != nil {
			global.SugaredLogger.Panicln(err)
		}
	}
}

// UploadSingleFile
// 上传单个文件专用，经过uuid加持后返回uuid和错误信息。
// 第二个入参为前端的关键词名称。
func UploadSingleFile(c *gin.Context, key string) (uniqueFilename string, err error) {
	_, header, err := c.Request.FormFile(key)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", err
	}
	if header.Size > global.Config.MaxSizeForUpload {
		return "", errors.New("文件过大")
	}
	id := uuid.NewString()
	header.Filename = id + "--" + header.Filename
	err = c.SaveUploadedFile(header, global.Config.FullPath+header.Filename)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", err
	}
	return id, nil
}

// UploadMultipleFiles 上传多个文件专用
// 经过uuid加持后返回唯一文件名和错误信息。
// 第二个入参为前端的关键词名称。
func UploadMultipleFiles(c *gin.Context, key string) (uuids []string, err error) {
	form, err := c.MultipartForm()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return nil, err
	}
	files := form.File[key]

	for _, file := range files {
		if file.Size > global.Config.MaxSizeForUpload {
			return nil, errors.New("文件过大")
		}
	}

	for _, file := range files {
		id := uuid.NewString()
		file.Filename = id + "--" + file.Filename
		err = c.SaveUploadedFile(file, global.Config.FullPath+file.Filename)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return nil, err
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}
