package handler

import (
	"net/http"
)

// 拦截器 验证请求合法性
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !isTokenValid(token) {
				http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
				return
			}
			h(w, r)
		})
}
