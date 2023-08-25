package db

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:   0,
	})

	ctx := context.Background()
	err := Redis.Ping(ctx).Err()
	if err != nil {
		log.Errorln("RedisConectionError: ", err.Error())
		return
	}
	log.Debugln("connected to redis")
}
