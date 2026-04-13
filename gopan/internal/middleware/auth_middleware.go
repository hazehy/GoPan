// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"errors"
	"gopan/gopan/helper"
	"gopan/gopan/models"
	"net/http"
	"strconv"

	"xorm.io/xorm"
)

type AuthMiddleware struct {
	engine *xorm.Engine
}

func NewAuthMiddleware(engine *xorm.Engine) *AuthMiddleware {
	return &AuthMiddleware{engine: engine}
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

		if err := m.ensureUserAvailable(uc.Identity); err != nil {
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

func (m *AuthMiddleware) ensureUserAvailable(identity string) error {
	if m.engine == nil {
		return errors.New("服务未就绪")
	}

	user := new(models.User)
	has, err := m.engine.Where("identity = ?", identity).Get(user)
	if err != nil {
		return errors.New("用户状态校验失败")
	}
	if !has || !user.DeletedAt.IsZero() {
		return errors.New("用户不存在")
	}
	if user.Status != 1 {
		return errors.New("用户已被禁用")
	}

	return nil
}
