package redis

import (
	"context"
	"github.com/guneyin/bist-tools/pkg/config"
	"github.com/guneyin/bist-tools/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type cache struct {
	c *redis.Client
	//l *logger.Logger
}

var rdb *cache

const ttl = time.Hour * 24

func Init() error {
	opt, err := redis.ParseURL(config.Cfg.Redis.URL)
	if err != nil {
		panic(err)
	}

	c := redis.NewClient(opt)

	_, err = c.Ping(context.TODO()).Result()
	if err != nil {
		logger.Error("redis: connection error \n%s", err.Error())

		return err
	}

	rdb = &cache{
		c: c,
		//l: logger.Log,
	}

	logger.Info("redis: connection successful")

	return nil
}

func Delete(ctx context.Context, k string) error {
	return rdb.c.Del(ctx, k).Err()
}

func Set(ctx context.Context, k string, v []byte) error {
	return rdb.c.Set(ctx, k, v, ttl).Err()
}

func Get(ctx context.Context, k string) ([]byte, error) {
	return rdb.c.Get(ctx, k).Bytes()
}

func SetH(ctx context.Context, k string, v any) error {
	err := rdb.c.HSet(ctx, k, v).Err()
	if err != nil {
		return err
	}

	return rdb.c.Expire(ctx, k, ttl).Err()
}

func GetH(ctx context.Context, k, f string) ([]byte, error) {
	return rdb.c.HGet(ctx, k, f).Bytes()
}

func GetHAll(ctx context.Context, k string) (map[string]string, error) {
	return rdb.c.HGetAll(ctx, k).Result()
}
