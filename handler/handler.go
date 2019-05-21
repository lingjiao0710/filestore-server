package handler

import (
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"time"
	"encoding/json"
	"github.com/lingjiao0710/filestore-server/meta"
	"github.com/lingjiao0710/filestore-server/util"
)

//UploadHandler: 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		//返回上传HTML页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w, "internel server error")
			return 
		}
		io.WriteString(w, string(data))
	}else if r.Method == "POST"{
		//接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Printf("get data failed, err: %s\n", err.Error())
			return 
		}

		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//创建本地文件
		newfile, err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Printf("creat file failed, err:%s\n", err.Error())
			return
		}
		defer newfile.Close()

		//复制数据
		fileMeta.Filesize, err = io.Copy(newfile, file)
		if err != nil {
			fmt.Printf("save data failed, err:%s\n", err.Error())
			return 
		}

		//移动文件指针到文件首
		newfile.Seek(0, 0)
		//计算文件SHA1
		fileMeta.FileSha1 = util.FileSha1(newfile)
		meta.UpdateFileMeta(fileMeta)
		//重定向到suc路由
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

//UploadSucHandler: 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "Upload success!")
}

//GetFileMetaHandler: 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	//filehash := r.Form["filehash"][0]
	filehash := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)

}

func DownloadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	fmt.Printf("fsha1:%s\n", fsha1)
	fMeta := meta.GetFileMeta(fsha1)

	f, err := os.Open(fMeta.Location)
	if err != nil {
		errout := fmt.Sprintf("open %s failed", fMeta.Location)
		w.Write([]byte(errout))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fMeta.FileName+"\"")
	w.Write(data)
}


func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return 
	}

	if r.Method != "POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return 
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//FileDeleteHandler： 删除文件接口
func FileDeleteHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")


	//删除文件
	fMeta := meta.GetFileMeta(fileSha1)
		err := os.Remove(fMeta.Location)
	if err != nil {
		fmt.Printf("remove %s failed, err:%s\n", fMeta.Location, err.Error())
		return 
	}
	
	//删除元信息
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}