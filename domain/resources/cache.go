package resources

import (
	"context"
	"time"

	"github.com/Mestrace/orderbook/conf"
	"github.com/go-redis/redis/v9"
)

var (
	redisClient *redis.Client
)

func InitRedis() error {
	var (
		rconf = conf.Get().Redis
	)
	redisClient = redis.NewClient(&redis.Options{
		Addr:         rconf.Addr,
		Password:     rconf.Password,
		DB:           rconf.Db,
		DialTimeout:  time.Duration(rconf.DialTimeout),
		ReadTimeout:  time.Duration(rconf.ReadTimeout),
		WriteTimeout: time.Duration(rconf.WriteTimeout),
	})

	{
		cmd := redisClient.Ping(context.TODO())
		if err := cmd.Err(); err != nil {
			return err
		}
	}

	return nil
}

func GetRedisClient() *redis.Client {
	return redisClient
}
