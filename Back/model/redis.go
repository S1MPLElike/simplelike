package model

import (
	"GoWorkingProject/Back/config"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// 全局Redis客户端和上下文（所有Redis操作都用这个）
var (
	RedisClient *redis.Client
	Ctx         = context.Background() // Redis操作上下文
)

// InitRedis 初始化Redis连接（只执行一次）
func InitRedis() error {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Addr,     // Redis地址+端口
		Password: config.RedisConfig.Password, // Redis密码
		DB:       config.RedisConfig.Database, // 选择数据库
	})

	// 测试Redis连接是否成功
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	fmt.Println("✅ Redis连接成功")
	return nil
}

// 设置Redis键值对（带过期时间）
func RedisSet(key string, value string, expire time.Duration) error {
	return RedisClient.Set(Ctx, key, value, expire).Err()
}

// 获取Redis键值对
func RedisGet(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

// 删除Redis键
func RedisDel(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}

// RedisSetNX 设置Redis键值对（带过期时间），如果键不存在则设置成功
func RedisSetNX(key string, value string, expire time.Duration) (bool, error) {
	return RedisClient.SetNX(Ctx, key, value, expire).Result()
}
