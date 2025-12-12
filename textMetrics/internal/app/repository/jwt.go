package repository

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

const jwtPrefix = "jwt:"

func (c *Redis) getJwtKey(token string) string {
	logrus.Info(token)
	return c.cfg.RedisAppPrefix + jwtPrefix + token
}

func (c *Redis) SetBlackListJWT(ctx context.Context, token string, jwtTTL time.Duration) error {
	return c.client.Set(ctx, c.getJwtKey(token), true, jwtTTL).Err()
}

func (c *Redis) GetBlackListJWT(ctx context.Context, token string) error {
	return c.client.Get(ctx, c.getJwtKey(token)).Err()
}
