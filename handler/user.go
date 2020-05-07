package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/JJMeg/filestore_server/db"
	"github.com/JJMeg/filestore_server/util"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler: 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	}
	r.ParseForm()

	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))

	suc := db.UserSignUp(username, enc_passwd)
	if suc {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("fail"))
	}
}

// SignInHandler：登陆接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	//	1. 校验用户名和密码
	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	pwdChecked := db.UserSignIn(username, enc_passwd)
	if !pwdChecked {
		w.Write([]byte("fail"))
		return
	}

	//  2. 生成访问凭证 token
	token := genToken(username)
	// token写入数据库
	upRes := db.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("fail"))
		return
	}

	//  3. 登陆成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
	return
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	//fmt.Println("token: ", token)

	// 验证token
	//if !isTokenValid(token) {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}

	//	查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func isTokenValid(token string) bool {
	//判断token的时效性，是否过期
	//从数据库查到token，对比是否一致
	if len(token) != 40 {
		return false
	}
	return true
}
func genToken(username string) string {
	//	40位字符：md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
