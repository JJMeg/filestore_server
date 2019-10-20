package meta

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init()  {
	fileMetas = make(map[string]FileMeta)
}

//添加fileMeta内的元素
func UpdateFileMeta(fmeta FileMeta)  {
	fileMetas[fmeta.FileSha1] = fmeta
}

//获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta  {
	return fileMetas[fileSha1]
}

