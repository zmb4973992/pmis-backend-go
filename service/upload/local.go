package upload

import (
	"errors"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"strconv"
)

type Local struct{}

func (l *Local) UploadSingleFile(fileHeader *multipart.FileHeader) (fileName string, err error) {
	if fileHeader.Size > global.Config.MaxSize {
		return "", errors.New("文件过大")
	}

	storagePath := global.Config.UploadConfig.StoragePath

	file := model.File{
		Name: fileHeader.Filename,
		Mode: "local",
		Path: storagePath,
		Size: int(fileHeader.Size >> 20), //MB
	}

	err = global.DB.Create(&file).Error
	if err != nil {
		return "", err
	}

	// 给文件名添加id和时间
	fileName = strconv.FormatInt(file.ID, 10) + "--" + fileHeader.Filename
	err = saveUploadedFile(fileHeader, storagePath+fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (l *Local) UploadMultipleFiles(fileHeaders []*multipart.FileHeader) (fileNames []string, err error) {
	for i := range fileHeaders {
		if fileHeaders[i].Size > global.Config.UploadConfig.MaxSize {
			return nil, errors.New("文件过大")
		}
	}

	storagePath := global.Config.UploadConfig.StoragePath

	for i := range fileHeaders {
		//保存记录
		var record model.File
		record.Name = fileHeaders[i].Filename
		record.Path = storagePath
		record.Mode = "local"
		record.Size = int(fileHeaders[i].Size) >> 20 // MB
		err = global.DB.Create(&record).Error
		if err != nil {
			return nil, err
		}

		fileName := strconv.FormatInt(record.ID, 10) + "--" + fileHeaders[i].Filename
		fileNames = append(fileNames, fileName)

		err = saveUploadedFile(fileHeaders[i], storagePath+fileName)
		if err != nil {
			return nil, err
		}
	}
	return fileNames, nil
}

func (l *Local) Delete(id int64) error {
	if id == 0 {
		return nil
	}

	var record model.File
	err := global.DB.Where(&model.File{BasicModel: model.BasicModel{ID: id}}).
		First(&record).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if record.ID > 0 {
		filePath := record.Path + strconv.FormatInt(id, 10)
		fileName := record.Name
		_ = os.Remove(filePath + "--" + fileName)

		err = global.DB.Delete(&record).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
		}
		return err
	}
	return nil
}

// 仿照gin c.SaveUploadedFile的写法
func saveUploadedFile(fileHeader *multipart.FileHeader, destination string) error {
	//打开、读取文件
	openedFile, err := fileHeader.Open()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}
	defer openedFile.Close()

	//创建空的新文件
	createdFile, err := os.Create(destination)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}
	defer createdFile.Close()

	//把打开的文件内容复制到新文件中
	_, err = io.Copy(createdFile, openedFile)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	return nil
}
