package db

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var Redis *redis.Client

func InitRedis() {
	var op *redis.Options
	if os.Getenv("ENV") == "dev" {
		op = &redis.Options{
			Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			DB:   0,
		}
	} else if os.Getenv("ENV") == "prod" {
		op = &redis.Options{
			Addr:      fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password:  os.Getenv("REDIS_PASSWORD"),
			TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
			DB:        0,
		}
	}
	Redis = redis.NewClient(op)

	ctx := context.Background()
	err := Redis.Ping(ctx).Err()
	if err != nil {
		log.Errorln("RedisConectionError: ", err.Error())
		return
	}
	log.Debugln("connected to redis")
}
