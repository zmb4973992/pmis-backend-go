package upload

import (
	"errors"
	"math"
	"mime/multipart"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"strconv"
)

type Local struct {
}

func (l *Local) UploadMultipleFiles(userId int64, fileHeaders []*multipart.FileHeader) (fileNames []string, err error) {
	for i := range fileHeaders {
		if fileHeaders[i].Size > global.Config.UploadConfig.MaxSize {
			return nil, errors.New("文件过大")
		}
	}

	//storagePath := global.Config.UploadConfig.StoragePath

	for i := range fileHeaders {
		//保存记录
		var record model.File
		record.Creator = &userId
		record.Name = fileHeaders[i].Filename

		record.SizeMB = math.Round(float64(fileHeaders[i].Size)/(1024*1024)*100) / 100
		err = global.DB.Create(&record).Error
		if err != nil {
			return nil, err
		}

		fileName := strconv.FormatInt(record.Id, 10) + "--" + fileHeaders[i].Filename
		fileNames = append(fileNames, fileName)

		//err = saveUploadedFile(fileHeaders[i], storagePath+fileName)
		//if err != nil {
		//	return nil, err
		//}
	}
	return fileNames, nil
}
