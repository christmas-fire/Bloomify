package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func AddTokenToRedis(ctx context.Context, client *redis.Client, username string, token string) error {
	err := client.Set(ctx, username, token, 0).Err()
	if err != nil {
		return fmt.Errorf("error insert JWT into Redis: %w", err)
	}

	log.Printf("token successfully added for user '%s'\n", username)

	return nil
}
