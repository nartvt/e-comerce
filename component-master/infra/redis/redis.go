package redis

import (
	"component-master/config"
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func InitRedis(rd *config.RedisConfig) (*redis.Client, error) {
	opt, err := redis.ParseURL(rd.BuildRedisConnectionString())
	if err != nil {
		return nil, err
	}
	opt.PoolSize = rd.MaxIdle
	opt.DialTimeout = rd.DialTimeout
	opt.ReadTimeout = rd.ReadTimeout
	opt.WriteTimeout = rd.WriteTimeout
	opt.Password = rd.Password
	opt.DB = rd.DB

	client := redis.NewClient(opt)
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	slog.Info(fmt.Sprintf("Ping to redis success: %s", pong))
	return client, nil
}
