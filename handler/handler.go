package handler

import (
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"fmt"
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

		//创建本地文件
		newfile, err := os.Create("./" + head.Filename)
		if err != nil{
			fmt.Printf("creat file failed, err:%s\n", err.Error())
			return
		}
		defer newfile.Close()

		//复制数据
		_, err = io.Copy(newfile, file)
		if err != nil {
			fmt.Printf("save data failed, err:%s\n", err.Error())
			return 
		}

		//重定向到suc路由
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

//UploadSucHandler: 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "Upload success!")
}

