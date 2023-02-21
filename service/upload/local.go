package upload

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

type Local struct{}

func (l *Local) UploadSingleFile(fileHeader *multipart.FileHeader) (accessPath string, fileName string, err error) {
	if fileHeader.Size > global.Config.MaxSize {
		return "", "", errors.New("文件过大")
	}
	// 给文件名添加uuid和时间，确保唯一性
	id := uuid.NewString()
	//formattedTime := time.Now().Format("2006-01-02-15-04-05")
	fileName = id + "--" + fileHeader.Filename
	storagePath := global.Config.UploadConfig.StoragePath
	accessPath = global.Config.DownloadConfig.AccessPath
	err = saveUploadedFile(fileHeader, storagePath+fileName)
	if err != nil {
		return "", "", err
	}
	return accessPath, fileName, nil
}

func (l *Local) UploadMultipleFiles(fileHeaders []*multipart.FileHeader) (accessPath string, fileNames []string, err error) {
	for i := range fileHeaders {
		if fileHeaders[i].Size > global.Config.UploadConfig.MaxSize {
			return "", nil, errors.New("文件过大")
		}
	}

	storagePath := global.Config.UploadConfig.StoragePath
	accessPath = global.Config.DownloadConfig.AccessPath

	for i := range fileHeaders {
		// 给文件名添加uuid和时间，确保唯一性
		id := uuid.NewString()
		//formattedTime := time.Now().Format("2006-01-02 15-04-05")
		fileName := id + "--" + fileHeaders[i].Filename
		err := saveUploadedFile(fileHeaders[i], storagePath+fileName)
		if err != nil {
			return "", nil, err
		}

		//保存记录
		var record model.File
		record.UUID = id
		record.InitialFileName = fileHeaders[i].Filename
		record.StoredFileName = fileName
		record.StoragePath = storagePath
		record.AccessPath = accessPath
		record.Size = int(fileHeaders[i].Size) >> 20 // MB
		global.DB.Create(&record)
		fileNames = append(fileNames, fileName)
	}
	return accessPath, fileNames, nil
}

func (l *Local) Delete(UUID string) error {
	var record model.File
	err := global.DB.Where("uuid = ?", UUID).First(&record).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if record.ID != 0 {
		storagePath := record.StoragePath
		storedFileName := record.StoredFileName
		_ = os.Remove(storagePath + storedFileName)

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
