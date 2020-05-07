package db

import (
	"fmt"
	"github.com/JJMeg/filestore_server/db/mysql"
)

// UserSignUp: 通过用户名及密码注册用户
func UserSignUp(username, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %v", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Printf("Failed to insert,err: %v", err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil && rf > 0 {
		return true
	}
	return false
}

// UserSignIn: 判断账号密码
func UserSignIn(username, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %v", err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Printf("Failed to select,err: %v", err.Error())
		return false
	} else if rows == nil {
		fmt.Printf("Not found,err: %v", err.Error())
		return false
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateToken: 刷新用户登陆的token
func UpdateToken(username, token string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %v", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, token)
	if err != nil {
		fmt.Printf("Failed to update token, err: %v", err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil && rf > 0 {
		return true
	}
	return false
}

type User struct {
	Username     string
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	SignupAt     string
	LastActiveAt string `json:"last_active_at"`
	Status       int    `json:"status"`
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %v", err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Printf("Failed to update token, err: %v", err.Error())
		return user, err
	}

	return user, nil
}
