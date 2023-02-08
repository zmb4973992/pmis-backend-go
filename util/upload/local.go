package upload

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"pmis-backend-go/global"
	"time"
)

type Local struct{}

func (l *Local) UploadSingleFile(fileHeader *multipart.FileHeader) (storagePath string, fileName string, err error) {
	if fileHeader.Size > global.Config.MaxSize {
		return "", "", errors.New("文件过大")
	}
	// 给文件名添加uuid和时间，确保唯一性
	id := uuid.NewString()
	formattedTime := time.Now().Format("2006-01-02-15-04-05")
	fileName = id + "--" + formattedTime + "--" + fileHeader.Filename
	// 拼接路径和文件名
	storagePath = global.Config.UploadConfig.Path
	// 读取文件
	openedFile, err := fileHeader.Open()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", "", err
	}
	defer openedFile.Close()
	// 创建文件
	createdFile, err := os.Create(storagePath + fileName)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", "", err
	}
	defer createdFile.Close()
	// 把文件内容复制到新生成的文件中
	_, err = io.Copy(createdFile, openedFile) // 传输（拷贝）文件
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return "", "", err
	}
	return storagePath, fileName, nil
}

func (l *Local) UploadMultipleFiles(fileHeaders []*multipart.FileHeader) (storagePath string, fileNames []string, err error) {
	for i := range fileHeaders {
		if fileHeaders[i].Size > global.Config.UploadConfig.MaxSize {
			return "", nil, errors.New("文件过大")
		}
	}

	storagePath = global.Config.UploadConfig.Path

	for i := range fileHeaders {
		// 给文件名添加uuid和时间，确保唯一性
		id := uuid.NewString()
		formattedTime := time.Now().Format("2006-01-02-15-04-05")
		fileName := id + "--" + formattedTime + "--" + fileHeaders[i].Filename
		fileNames = append(fileNames, fileName)
		// 读取文件
		openedFile, err := fileHeaders[i].Open()
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return "", nil, err
		} else {
			defer openedFile.Close()
		}
		// 创建文件
		createdFile, err := os.Create(storagePath + fileName)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return "", nil, err
		}
		defer createdFile.Close()
		// 把文件内容复制到新生成的文件中
		_, err = io.Copy(createdFile, openedFile) // 传输（拷贝）文件
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return "", nil, err
		}
	}
	return storagePath, fileNames, nil
}

func (l *Local) Delete(key string) error {
	return nil
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
