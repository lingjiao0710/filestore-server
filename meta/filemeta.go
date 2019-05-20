package meta

//FileMeta: 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	Filesize int64
	Location string
	UploadAt string
}


var fileMetas map[string]FileMeta

//初始化fileMetas
func init() {
	fileMetas = make(map[string]FileMeta)
}

//UpdateFileMeta: 新增、更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//GetFileMeta: 通过sha1获取文件元信息

func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}
