// backend/utils/aferos/localfile/localfile.go

package aferos

import (
	"fmt"
)

type LocalFile struct {
	*FileStore
}

// 新建本地文件
func NewLocalFile(basePath string) (*LocalFile, error) {
	fs, err := NewFileStore(basePath)
	if err != nil {
		return nil, err
	}
	return &LocalFile{fs}, nil
}

// 保存文件切片
func (lf *LocalFile) SaveFileSlice(fileHash, sliceIndex string, data []byte) error {
	// TODO：文件名字之前只是片的索引，目前修改为文件hash+片的索引
	sliceName := fmt.Sprintf("%s %s", fileHash, sliceIndex) // 切片名称(文件哈希、切片索引)
	return lf.Write(fileHash, sliceName, data)
}

// 删除basePath文件路径下的所有文件夹和文件
func (lf *LocalFile) ClearBasePath() error {
	err := lf.DeleteAll("")
	if err != nil {
		return err
	}
	return nil
}
