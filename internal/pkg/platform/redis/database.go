package redis

import (
	"context"
	"log"
	"time"
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/adapter/kvadapter"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedis() (kvadapter.RepoAdapter, error) {
	cfg := conf.GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return Redis{
		rdb: client,
	}, nil
}

func (r Redis) Delete(key string) (err error) {
	res := r.rdb.Del(context.Background(), key)
	err = res.Err()
	return
}

// SetValueUntilChannelClose will store value and until isOpen is false.
func (r Redis) SetValueUntilChannelClose(key string, data string, ttl int, isOpen *bool) {
	log.Printf("On key %s this value will set: %s", key, data)
	go func() {
		ticker := time.NewTicker(time.Duration(ttl-2) * time.Second)
		defer func() {
			ticker.Stop()
			r.Delete(key)
		}()
		for *isOpen {
			op1 := r.rdb.Set(context.Background(), key, data, time.Duration(ttl)*time.Second)
			log.Printf("On key %s this value was set: %s", key, data)
			if op1.Err() != nil {
				//Error
				log.Printf("Error when value %s is set", data)
				continue
			}
			<-ticker.C
		}
	}()
}

func (r Redis) GetValue(key string) (res string, err error) {
	op2 := r.rdb.Get(context.Background(), key)
	if err = op2.Err(); err != nil {
		return
	}
	res, err = op2.Result()
	return
}
