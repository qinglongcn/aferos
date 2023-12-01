// backend/utils/aferos/filestore/filestore.go
// 定义共享的基类和方法

package aferos

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type FileStore struct {
	Fs       afero.Fs
	BasePath string
}

// 新建文件存储
func NewFileStore(basePath string) (*FileStore, error) {
	fs := afero.NewOsFs()
	// MkdirAll 创建目录路径和所有尚不存在的父级。
	if err := fs.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &FileStore{
		Fs:       fs,
		BasePath: basePath,
	}, nil
}

// 创建文件
func (fs *FileStore) CreateFile(subDir, fileName string) error {
	filePath := filepath.Join(fs.BasePath, subDir, fileName)
	if exists, _ := afero.Exists(fs.Fs, filePath); !exists {
		// 文件不存在，创建文件
		if _, err := fs.Fs.Create(filePath); err != nil {
			return err
		}
	}
	return nil
}

// 写入文件
func (fs *FileStore) Write(subDir, fileName string, data []byte) error {
	if err := fs.Fs.MkdirAll(filepath.Join(fs.BasePath, subDir), 0755); err != nil {
		return err
	}
	filePath := filepath.Join(fs.BasePath, subDir, fileName)
	return afero.WriteFile(fs.Fs, filePath, data, 0644)
}

// 读取文件
func (fs *FileStore) Read(subDir, fileName string) ([]byte, error) {
	filePath := filepath.Join(fs.BasePath, subDir, fileName)
	if exists, _ := afero.Exists(fs.Fs, filePath); !exists {
		return nil, fmt.Errorf("文件未找到")
	}
	return afero.ReadFile(fs.Fs, filePath)
}

// 删除文件
func (fs *FileStore) Delete(subDir, fileName string) error {
	filePath := filepath.Join(fs.BasePath, subDir, fileName)
	return fs.Fs.Remove(filePath)
}

// 删除所有文件
func (fs *FileStore) DeleteAll(subDir string) error {
	filePath := filepath.Join(fs.BasePath, subDir)
	return fs.Fs.RemoveAll(filePath)
}

// 检查文件是否存在
func (fs *FileStore) Exists(subDir, fileName string) (bool, error) {
	filePath := filepath.Join(fs.BasePath, subDir, fileName)
	return afero.Exists(fs.Fs, filePath)
}

// 获取文件列表
func (fs *FileStore) ListFiles(subDir, partialName string) ([]string, error) {
	dirPath := filepath.Join(fs.BasePath, subDir)
	files, err := afero.ReadDir(fs.Fs, dirPath)
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, file := range files {
		if strings.Contains(file.Name(), partialName) {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList, nil
}

// CopyFile将一个文件从源路径复制到目标路径
func (fs *FileStore) CopyFile(srcFile, destFile, newFileName string) error {
	srcFilePath := filepath.Join(fs.BasePath, srcFile)
	destFilePath := destFile
	if err := fs.Fs.MkdirAll(destFilePath, 0755); err != nil {
		return err
	}
	src, err := fs.Fs.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer src.Close()
	name := filepath.Join(destFilePath, newFileName)
	dest, err := fs.Fs.Create(name)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return nil
}
