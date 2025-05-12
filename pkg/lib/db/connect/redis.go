package connect

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func RedisConnect(ctx context.Context, Host, Port, Password, Env string, dbnumber int) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", Host, Port)
	var rdb *redis.Client
	switch Env {
	case "development":
		rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379", // или "redis:6379" внутри Docker-сети
			DB:   0,
		})
	case "production":
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr, // или "redis:6379" внутри Docker-сети
			Password: Password,
			DB:       dbnumber,
		})
	}

	// Проверка подключения
	ping, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("Redis connect success:", ping)
	return rdb, nil
}
