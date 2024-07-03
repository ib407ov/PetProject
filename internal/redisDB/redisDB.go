package redisDB

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func ConnectRedisDB() *redis.Client {
	opt := &redis.Options{
		Addr:     viper.GetString("redis.uri"),
		Password: "",
		DB:       0,
	}

	client := redis.NewClient(opt)

	// Перевірка підключення
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	logrus.Println("Redis connect success")

	return client
}
