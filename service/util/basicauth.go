package util

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// 定义一个认证函数，检查用户名和密码
func authenticate1(username, password string) bool {
	// 这里应该从安全存储中获取正确的用户名和密码
	// 为了演示，我们使用硬编码的用户名和密码
	return username == "admin" && password == "het002402"
}

// 用于拦截请求并验证基本认证
func BasicAuth(handler http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从HTTP请求头中获取认证信息
		auth := r.Header.Get("Authorization")
		if auth == "" {
			// 如果没有提供认证信息，返回401 Unauthorized
			w.Header().Set("WWW-Authenticate", `Basic realm="Secure Area"`)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		// 解码认证信息
		decoded, err := base64.StdEncoding.DecodeString(auth[6:]) // 跳过"Basic "前缀
		if err != nil {
			http.Error(w, "Error decoding auth token", http.StatusBadRequest)
			return
		}

		// 将解码后的字符串分割为用户名和密码
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			http.Error(w, "Invalid auth token", http.StatusBadRequest)
			return
		}

		// 验证用户名和密码
		//if !authenticate(creds[0], creds[1]) {
		if username != creds[0] || password != creds[1] {
			w.Header().Set("WWW-Authenticate", `Basic realm="Secure Area"`)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		// 认证通过，调用原始的handler
		handler(w, r)
	}
}
