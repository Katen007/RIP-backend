package repository

import (
	"context"
	"fmt"
	"lab1_rip/internal/app/config"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	cfg    *config.Config
	client *redis.Client
}

func NewRedis(cfg *config.Config) (*Redis, error) {
	client := Redis{}

	client.cfg = cfg

	redisClient := redis.NewClient(&redis.Options{
		Password:    cfg.RedisPassword,
		Username:    cfg.RedisUser,
		Addr:        cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
		DB:          0,
		DialTimeout: cfg.RedisDialTimeoutSec,
		ReadTimeout: cfg.RedisReadTimeoutSec,
	})

	client.client = redisClient
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return &client, nil
}

func (c *Redis) Close() error {
	return c.client.Close()
}
