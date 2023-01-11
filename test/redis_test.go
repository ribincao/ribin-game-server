package test

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func TestRedisConn(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func TestRedisSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Set("ping", "pong", 0).Err()
	if err != nil {
		panic(err)
	}
}

func TestRedisGet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := client.Get("ping").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}