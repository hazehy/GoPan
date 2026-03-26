package models

import (
	"gopan/gopan/internal/config"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func Init(MysqlDatabase string) *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", MysqlDatabase)
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
	}
	return engine
}

func InitRedis(c config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: "",
		DB:       0,
	})
	return rdb
}
