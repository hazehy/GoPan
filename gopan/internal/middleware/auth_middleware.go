// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"gopan/gopan/helper"
	"net/http"
	"strconv"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	// 从请求头中获取Authorization字段，解析JWT token，验证用户身份，并将用户信息添加到请求头中
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		uc, err := helper.AnalyzeToken(auth)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		r.Header.Set("UserId", string(rune(uc.Id)))
		r.Header.Set("UserIdentity", uc.Identity)
		r.Header.Set("UserName", uc.Name)
		r.Header.Set("UserRole", strconv.Itoa(uc.Role))
		next(w, r)
	}
}
