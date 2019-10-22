package mysql

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func init() {
	db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to mysql, err: %v", err)
		os.Exit(1)
	}
}

//返回数据库对象
func DBConn() *sql.DB {
	return db
}
