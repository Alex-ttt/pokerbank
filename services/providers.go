package services

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

var Db *sql.DB

var RedisClient redis.Conn

func SetRedisKey(redisClient redis.Conn, key string, value string) error {
	_, err := redisClient.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFromRedis(redisClient redis.Conn, key string) error {
	_, err := redisClient.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func GetFromRedis(redisClient redis.Conn, key string) (string, error) {
	result, err := redis.String(redisClient.Do("GET", key))
	if err != nil {
		return "", err
	}

	return result, nil
}

func SetRedisKeyWithExpiration(redisClient redis.Conn, key string, value string, expiration time.Duration) error {
	err := SetRedisKey(redisClient, key, value)
	if err != nil {
		return err
	}

	seconds := int(expiration / time.Second)
	_, err = redisClient.Do("EXPIRE", key, seconds)
	if err != nil {
		return err
	}

	return nil
}

func InitRedis() {
	redisUri := os.Getenv("REDIS_URL")
	var redisClient redis.Conn
	var err error
	if len(redisUri) > 0 {
		redisClient, err = redis.DialURL(redisUri)
		if err != nil {
			panic(err)
		}
	} else {
		redisClient, err = redis.Dial("tcp", ":6379")
		if err != nil {
			panic(err)
		}
	}

	RedisClient = redisClient
	//
	//
	////Initializing redis
	//dsn := os.Getenv("REDIS_DSN")
	//if len(dsn) == 0 {
	//	dsn = "localhost:6379"
	//}
	//RedisClient = redis.NewClient(&redis.Options{
	//	Addr: dsn, //redis port
	//})
	//_, err := RedisClient.Ping().Result()
	//if err != nil {
	//	panic(err)
	//}
}
