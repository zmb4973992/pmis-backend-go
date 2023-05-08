package upload

import (
	"errors"
	"github.com/yitter/idgenerator-go/idgen"
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
	// 给文件名添加snowID和时间
	snowID := idgen.NextId()
	fileName = strconv.FormatInt(snowID, 10) + "--" + fileHeader.Filename
	storagePath := global.Config.UploadConfig.StoragePath
	err = saveUploadedFile(fileHeader, storagePath+fileName)
	if err != nil {
		return "", err
	}

	file := model.File{
		BasicModel: model.BasicModel{
			SnowID: snowID,
		},
		Name: fileHeader.Filename,
		Mode: "local",
		Path: storagePath,
		Size: int(fileHeader.Size >> 20), //MB
	}

	err = global.DB.Create(&file).Error
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
		// 给文件名添加snowID和时间
		snowID := idgen.NextId()
		//formattedTime := time.Now().Format("2006-01-02 15-04-05")
		fileName := strconv.FormatInt(snowID, 10) + "--" + fileHeaders[i].Filename
		err := saveUploadedFile(fileHeaders[i], storagePath+fileName)
		if err != nil {
			return nil, err
		}

		//保存记录

		var record model.File
		record.SnowID = snowID
		record.Name = fileHeaders[i].Filename
		record.Path = storagePath
		record.Size = int(fileHeaders[i].Size) >> 20 // MB
		global.DB.Create(&record)
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil
}

func (l *Local) Delete(snowID int64) error {
	if snowID == 0 {
		return nil
	}

	var record model.File
	err := global.DB.Where(&model.File{BasicModel: model.BasicModel{SnowID: snowID}}).
		First(&record).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if record.SnowID > 0 {
		filePath := record.Path + strconv.FormatInt(snowID, 10)
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
