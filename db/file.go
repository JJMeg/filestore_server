package db

import (
	"database/sql"
	"filestore_server/db/mysql"
	"fmt"
)

//文件上传完成，保存meta新纪录
func OnFileUploadFinished(
	filehash string,
	filename string,
	filesize int64,
	fileaddr string) bool {

	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`, `file_size`, `file_addr`, `status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Printf("Failed to prepare statement,err:", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			//	未产生新表记录
			fmt.Printf("File with hash:%s has benn uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 从mysql获取元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	defer stmt.Close()

	tfile := TableFile{}

	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return &tfile, nil
}
