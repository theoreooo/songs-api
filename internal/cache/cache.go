package cache

import (
	"context"
	"os"
	"songs/internal/logger"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		logger.Log.Error("REDIS_ADDR не задан")
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	if err := Rdb.Ping(ctx).Err(); err != nil {
		logger.Log.Errorf("Не удалось подключиться к Redis: %v", err)
		return
	}
	logger.Log.Info("Успешно подключились к Redis")
}
