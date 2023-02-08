package upload

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/util"
)

// Oss 对象存储接口
type Oss interface {
	UploadSingleFile(fileHeader *multipart.FileHeader) (storagePath string, fileName string, err error)
	UploadMultipleFiles(fileHeaders []*multipart.FileHeader) (storagePath string, fileNames []string, err error)
	Delete(key string) error
}

func NewOss() Oss {
	switch global.Config.OssConfig.Type {
	case "local":
		return &Local{}
	default:
		return &Local{}
	}
}

func Init() {
	//检查上传文件的文件夹是否存在
	exists := util.PathExistsOrNot(global.Config.Path)
	//如果不存在就创建
	if !exists {
		err := os.MkdirAll(global.Config.Path, os.ModePerm)
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
	if header.Size > global.Config.MaxSize {
		return "", errors.New("文件过大")
	}
	id := uuid.NewString()
	header.Filename = id + "--" + header.Filename
	err = c.SaveUploadedFile(header, global.Config.Path+header.Filename)
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
		if file.Size > global.Config.MaxSize {
			return nil, errors.New("文件过大")
		}
	}

	for _, file := range files {
		id := uuid.NewString()
		file.Filename = id + "--" + file.Filename
		err = c.SaveUploadedFile(file, global.Config.Path+file.Filename)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return nil, err
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}
