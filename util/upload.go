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
	res := PathExistsOrNot(global.Config.FullPath)
	//如果不存在就创建
	if res == false {
		err := os.MkdirAll(global.Config.FullPath, os.ModePerm)
		if err != nil {
			global.SugaredLogger.Panicln(err)
		}
	}
}

// UploadSingleFile
// 上传单个文件专用，经过uuid加持后返回唯一文件名和错误信息。
// 第二个入参为前端的关键词名称。
func UploadSingleFile(c *gin.Context, key string) (uniqueFilename string, err error) {
	file, err := c.FormFile(key)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", err
	}
	if file.Size > global.Config.MaxSizeForUpload {
		return "", errors.New("文件过大")
	}
	id := uuid.New().String()
	file.Filename = id + "--" + file.Filename
	err = c.SaveUploadedFile(file, global.Config.FullPath+file.Filename)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", err
	}
	return file.Filename, nil
}

// UploadMultipleFiles 上传多个文件专用，经过uuid加持后返回唯一文件名和错误信息，文件名之间用竖线 | 分隔。
// 第二个入参为前端的关键词名称。
func UploadMultipleFiles(c *gin.Context, key string) (uniqueFilenames []string, err error) {
	form, err := c.MultipartForm()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return nil, err
	}
	files := form.File[key]
	var fileNames []string

	for _, file := range files {
		if file.Size > global.Config.MaxSizeForUpload {
			return nil, errors.New("文件过大")
		}
	}

	for _, file := range files {
		id := uuid.New().String()
		file.Filename = id + "--" + file.Filename
		err = c.SaveUploadedFile(file, global.Config.FullPath+file.Filename)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return nil, err
		}
		fileNames = append(fileNames, file.Filename)
	}
	//res := strings.Join(fileNames, "|")
	return fileNames, nil
}
