package db

import (
	"fmt"
	"time"

	"github.com/JJMeg/filestore_server/db/mysql"
)

//用户文件表结构
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// OnUserFileUploadFinished: 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`,`upload_at`) values (?,?,?,?,?)")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %v", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		fmt.Printf("Failed to insert,err: %v", err.Error())
		return false
	}

	if _, err := ret.RowsAffected(); err != nil {
		return false
	}
	return true
}

func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %v", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Printf("Failed to select,err: %v", err.Error())
		return nil, err
	}

	fmt.Println("rows: ", rows)
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileHash, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			break
		}
		userFiles = append(userFiles, ufile)
	}

	return userFiles, nil
}
