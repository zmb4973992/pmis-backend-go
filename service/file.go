package service

import (
	"errors"
	"io"
	"math"
	"mime/multipart"
	"os"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type FileGet struct {
	Id int64
}

type FileCreate struct {
	UserId     int64
	FileHeader *multipart.FileHeader
}

type FileDelete struct {
	Id int64
}

//以下为出参

type FileOutput struct {
	Creator      *int64  `json:"creator"`
	LastModifier *int64  `json:"last_modifier"`
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Url          string  `json:"url" gorm:"-"`
	SizeMB       float64 `json:"size_mb"`
	CreatedAt    *string `json:"created_at"`
}

func (f *FileGet) Get() (filePath string, fileName string, existed bool) {
	var record FileOutput
	err := global.DB.Model(model.File{}).
		Where("id = ?", f.Id).
		First(&record).Error
	if err != nil {
		return "", "", false
	}

	storagePath := global.Config.Upload.StoragePath
	filePath = storagePath + strconv.FormatInt(record.Id, 10) +
		"--" + record.Name
	//看该文件是否存在于服务器的文件夹中
	_, err = os.Stat(filePath)
	if err != nil {
		return "", "", false
	}

	return filePath, record.Name, true
}

func (f *FileCreate) Create() (fileId int64, url string, err error) {
	if f.FileHeader.Size > global.Config.Upload.MaxSize {
		return 0, "", errors.New("文件过大")
	}

	storagePath := global.Config.Upload.StoragePath

	file := model.File{
		BasicModel: model.BasicModel{
			Creator: &f.UserId,
		},
		Name:   f.FileHeader.Filename,
		SizeMB: math.Round(float64(f.FileHeader.Size)/(1024*1024)*100) / 100,
	}

	err = global.DB.Create(&file).Error
	if err != nil {
		return 0, "", err
	}

	// 给文件名添加id
	fileNameWithId := strconv.FormatInt(file.Id, 10) + "--" + f.FileHeader.Filename
	err = saveUploadedFile(f.FileHeader, storagePath+fileNameWithId)
	if err != nil {
		return 0, "", err
	}

	url = global.Config.Download.FullPath + strconv.FormatInt(file.Id, 10)

	return file.Id, url, nil
}

func (f *FileDelete) Delete() (errCode int) {
	var record model.File
	err := global.DB.Where("id = ?", f.Id).
		First(&record).Error
	if err != nil {
		return util.Success
	}

	storagePath := global.Config.Upload.StoragePath
	filePath := storagePath + strconv.FormatInt(record.Id, 10) +
		"--" + record.Name
	err = os.Remove(filePath)

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	err = global.DB.Where("id = ?", f.Id).Delete(&record).Error
	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	return util.Success
}

// 仿照gin c.SaveUploadedFile的写法
func saveUploadedFile(fileHeader *multipart.FileHeader, destination string) error {
	//打开、读取文件
	openedFile, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer openedFile.Close()

	//创建空的新文件
	createdFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer createdFile.Close()

	//把打开的文件内容复制到新文件中
	_, err = io.Copy(createdFile, openedFile)
	if err != nil {
		return err
	}

	return nil
}
