package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

// 用于向用户返回数据的对象，用于接收用户请求的请求指针
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//	返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
		io.WriteString(w,"Internal server error")
		}
		io.WriteString(w,string(data))

	} else if r.Method == "POST" {
		//	接收文件流及存储到本地目录
	}
}
