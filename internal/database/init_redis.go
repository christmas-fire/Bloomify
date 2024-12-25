package database

import (
	"context"
	"log"

	"github.com/christmas-fire/Bloomify/configs"
	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context) *redis.Client {
	cfg, err := configs.LoadConfigRedis("./configs")
	if err != nil {
		log.Fatal(err)
	}

	addr := cfg.Addr
	password := cfg.Password
	db := cfg.DB

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("error connect to Redis: %v", err)
	}

	log.Println("success connect to Redis")

	return client
}
