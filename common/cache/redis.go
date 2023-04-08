package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type MODE uint

var rdb *redis.Client

const (
	Solo MODE = iota
	Cluster
	Invalid
)

// InitRedis returns a solo mode redis client
func InitRedis() (*redis.Client, error) {
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

	_, err := rdb.Ping(ctx).Result()

	return rdb, err
}

// InitRedisCluster returns a cluster mode redis client
func InitRedisCluster() (*redis.ClusterClient, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        viper.GetStringSlice("redis.clusterConfig.addrs"),
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()

	return rdb, err
}

func GetClient() *redis.Client {
	return rdb
}
