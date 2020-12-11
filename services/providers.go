package services

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	"os"
)

var Db *sql.DB

var RedisClient *redis.Client

func InitRedis() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}
