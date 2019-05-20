package main

import (
	"fmt"
	"net/http"
	"github.com/lingjiao0710/filestore-server/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Start server failed, err: %s\n", err.Error())
		return 
	}
}