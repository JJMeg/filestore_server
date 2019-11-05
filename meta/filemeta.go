package meta

import (
	"filestore_server/db"
	"fmt"
)

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

//添加fileMeta内的元素
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// 新增meta到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	success := db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
	fmt.Printf("upload to db success: ", success)
	return success
}

//获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//删除文件元信息对象
func RemoveFileMeta(filesha1 string) {
	delete(fileMetas, filesha1)
}
