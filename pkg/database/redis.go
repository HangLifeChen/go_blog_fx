package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"go_blog/pkg/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis initializes the redis connection
func NewRedis(cfg *config.Config) *redis.Client {
	conf := cfg.Database.Redis

	opt := &redis.Options{
		Network:         "tcp",
		Addr:            fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		DB:              conf.Db,
		MinIdleConns:    16,
		ConnMaxLifetime: 300 * time.Second,
		Password:        conf.Password,
	}
	if cfg.Database.Redis.EnableTLS {
		opt.TLSConfig = &tls.Config{}
	}
	rdb := redis.NewClient(opt)
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}
