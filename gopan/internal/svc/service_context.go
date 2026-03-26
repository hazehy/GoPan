// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"gopan/gopan/internal/config"
	"gopan/gopan/internal/middleware"
	"gopan/gopan/models"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config config.Config
	Engine *xorm.Engine
	RDB    *redis.Client
	Auth   rest.Middleware
	Admin  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Engine: models.Init(c.Mysql.Database),
		RDB:    models.InitRedis(c),
		Auth:   middleware.NewAuthMiddleware().Handle,
		Admin:  middleware.NewAdminMiddleware().Handle,
	}
}
