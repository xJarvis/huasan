package redis

import (
	redis "github.com/go-redis/redis/v8"
	"github.com/xjarvis/huashan/app/config"
	"github.com/xjarvis/huashan/log/logger"
)

var (
	Open bool
	Redis *redis.Client
)

func Initialize() {
	if Open && Redis != nil {
		return
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Get("redis.host"),
		Password: config.Get("redis.password"),
		//DB:       config.Get("redis.database"),
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	// todo Context使用未进行验证
	_,err := client.Ping(client.Context()).Result()
	if err == nil {
		Open = true
		Redis = client
		logger.Info("Redis init finish")
	} else {
		Open = false
		logger.Warn("Redis init failed")
	}
}

func UnInitialize() {
	if Redis != nil {
		Redis = nil
	}
	logger.Info("Redis release finish")
}