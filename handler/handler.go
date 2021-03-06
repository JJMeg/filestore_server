package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/JJMeg/filestore_server/db"
	"github.com/JJMeg/filestore_server/meta"
	"github.com/JJMeg/filestore_server/util"
)

// TryFastUploadHandler: 尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// 1.解析请求参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 2. 从文件表中查询有相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "妙传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// 4. 上传过则姜文件信息写入用户文件表，返回成功
	if suc := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize)); suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	} else {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
		}
		w.Write(resp.JSONBytes())
		return
	}
}

// FileQueryHandler: 批量查询文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	userFiles, err := db.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		fmt.Println("db.QueryUserFileMetas: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 用于向用户返回数据的对象，用于接收用户请求的请求指针
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//	返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal server error")
		}
		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
		fmt.Println("Here into file upload")
		//	接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data: %v", err)
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 创建本地文件接收文件流
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("create new file failed: %v", err)
			return
		}
		defer newFile.Close()

		//把新文件流考到文件buffer去
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("io copy save file error: %v", err)
			return
		}

		//算文件的hash值
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMetaDB(fileMeta)

		//TODO: 更新用户文件表记录
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Println("here username:", username)
		suc := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Upload failed"))
			return
		}
	}
}

//显示上传成功信息
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload successfully")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fileMeta != nil {
		data, err := json.Marshal(fileMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	} else {
		w.Write([]byte(`{"code":-1,"msg":"no such file"}`))
	}
}

//文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(filehash)

	f, err := os.Open(fMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()
	//小文件可用ioutil，大文件使用流的方式，每次读一小部分数据给客户端再刷新缓存
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//加上头让浏览器当成文件进行下载
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fMeta.FileName+"\"")
	w.Write(data)
}

//更新元信息接口
func FileUpdateMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	file_hash := r.Form.Get("filehash")
	new_fname := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(file_hash)
	curFileMeta.FileName = new_fname
	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//删除文件
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	file_hash := r.Form.Get("filehash")
	meta.RemoveFileMeta(file_hash)

	//做真正的物理文件删除
	fMeta := meta.GetFileMeta(file_hash)

	os.Remove(fMeta.Location)
	w.WriteHeader(http.StatusOK)
}
