package service

import (
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type FileGet struct {
	ID int64
}

type FileCreate struct {
	UserID     int64
	FileHeader *multipart.FileHeader
}

type FileDelete struct {
	ID int64
}

//以下为出参

type FileOutput struct {
	Creator      *int64  `json:"creator"`
	LastModifier *int64  `json:"last_modifier"`
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	SizeMB       float64 `json:"size_mb"`
}

func (f *FileGet) Get() (filePath string, fileName string, existed bool) {
	var record FileOutput
	err := global.DB.Model(model.File{}).
		Where("id = ?", f.ID).
		First(&record).Error
	if err != nil {
		return "", "", false
	}

	storagePath := global.Config.UploadConfig.StoragePath
	filePath = storagePath + strconv.FormatInt(record.ID, 10) +
		"--" + record.Name
	//看该文件是否存在于服务器的文件夹中
	_, err = os.Stat(filePath)
	if err != nil {
		return "", "", false
	}

	return filePath, record.Name, true
}

func (f *FileCreate) Create() (id int64, err error) {
	fmt.Println(f.FileHeader.Size)
	if f.FileHeader.Size > global.Config.MaxSize {
		return 0, errors.New("文件过大")
	}

	storagePath := global.Config.UploadConfig.StoragePath

	file := model.File{
		BasicModel: model.BasicModel{
			Creator: &f.UserID,
		},
		Name:   f.FileHeader.Filename,
		SizeMB: math.Round(float64(f.FileHeader.Size)/(1024*1024)*100) / 100,
	}

	err = global.DB.Create(&file).Error
	if err != nil {
		return 0, err
	}

	// 给文件名添加id
	fileNameWithID := strconv.FormatInt(file.ID, 10) + "--" + f.FileHeader.Filename
	err = saveUploadedFile(f.FileHeader, storagePath+fileNameWithID)
	if err != nil {
		return 0, err
	}

	return file.ID, nil
}

func (f *FileDelete) Delete() response.Common {
	var record model.File
	err := global.DB.Where("id = ?", f.ID).
		First(&record).Error
	if err != nil {
		return response.Success()
	}

	storagePath := global.Config.UploadConfig.StoragePath
	filePath := storagePath + strconv.FormatInt(record.ID, 10) +
		"--" + record.Name
	err = os.Remove(filePath)

	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	err = global.DB.Where("id = ?", f.ID).Delete(&record).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	return response.Success()
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
