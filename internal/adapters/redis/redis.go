package redis

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr string) *goredis.Client {
	// Fail at startup rather than accepting requests with an unavailable database.
	rdb := goredis.NewClient(&goredis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Printf("connected to redis at %s", addr)

	return rdb
}
