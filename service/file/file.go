package file

import (
	"fmt"
	"io"
	"msm/config"
	"msm/log"
	"msm/model"
	"os"
)

type fileService struct{}

var FileService = new(fileService)

func (f *fileService) ReadFileFromPath(path string) (result []byte, err error) {
	fi, err := os.Open(path)
	if err != nil {
		return
	}
	defer fi.Close()
	fileInfo, err := fi.Stat()
	if err != nil {
		return
	}
	if size := float64(fileInfo.Size()) / 1e6; size > config.CF.FileSizeLimit {
		err = fmt.Errorf("写入数据大小%vMB,超过%vMB限制", size, config.CF.FileSizeLimit)
		return
	}
	result, err = io.ReadAll(fi)
	if err != nil {
		return
	}
	log.Logger.Debugw("文件写入成功", "path", path)
	return
}

func (f *fileService) UpdateFileData(filePath string, file io.Reader, size int64) error {
	if size := float64(size) / 1e6; size > config.CF.FileSizeLimit {
		return fmt.Errorf("写入数据大小%vMB,超过%vMB限制", size, config.CF.FileSizeLimit)
	}
	fi, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer fi.Close()
	if _, err = io.Copy(fi, file); err != nil {
		return err
	}
	log.Logger.Debugw("文件写入成功", "path", filePath)
	return nil
}

func (f *fileService) GetFileAndDirByPath(srcPath string) ([]model.FileStruct, error) {
	result := []model.FileStruct{}
	files, err := os.ReadDir(srcPath)
	if err != nil {
		return result, err
	}
	for _, file := range files {
		result = append(result, model.FileStruct{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		})
	}
	return result, nil
}

func (f *fileService) CreateNewDir(path string, name string) error {
	_, err := os.Create(path + name)
	return err
}

func (f *fileService) CreateNewFile(path string, name string) error {
	return os.MkdirAll(path+name, os.ModeDir)
}
