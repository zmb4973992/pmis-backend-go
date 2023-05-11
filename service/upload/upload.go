package upload

import (
	"mime/multipart"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/util"
)

// Oss 对象存储接口
type Oss interface {
	UploadSingleFile(fileHeader *multipart.FileHeader) (fileName string, err error)
	UploadMultipleFiles(fileHeaders []*multipart.FileHeader) (fileNames []string, err error)
	Delete(id int64) error
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
	exists := util.PathExistsOrNot(global.Config.StoragePath)
	//如果不存在就创建
	if !exists {
		err := os.MkdirAll(global.Config.StoragePath, os.ModePerm)
		if err != nil {
			global.SugaredLogger.Panicln(err)
		}
	}
}
