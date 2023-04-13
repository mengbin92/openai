package cache

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	rdb      *redis.Client
	initOnce sync.Once
)

// InitRedis returns a solo mode redis client
func Init() error {
	var err error

	initOnce.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:         viper.GetString("redis.addr"),
			DB:           viper.GetInt("redis.db"),
			Password:     viper.GetString("redis.password"),
			DialTimeout:  60 * time.Second,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		})

		ctx, cancal := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancal()

		_, err = rdb.Ping(ctx).Result()
	})

	return err
}

func Get() *redis.Client {
	if rdb == nil {
		panic("local redis is nil")
	}
	return rdb
}
