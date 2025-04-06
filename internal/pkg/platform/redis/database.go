package redis

import (
	"context"
	"fmt"
	"log"
	"time"
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/adapter/kvadapter"
	"websocket_client/internal/pkg/platform/zaplogger"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	rdb *redis.Client
	lgr zaplogger.Logger
}

// NewRedis initializes a new connection to the redis instance
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

// Delete deletes the redis key or returns an error
func (r Redis) Delete(key string) (err error) {
	res := r.rdb.Del(context.Background(), key)
	err = res.Err()
	return
}

// SetValueUntilChannelClose will continuously set a key value pair every ttl - some small number until isOpen is false.
func (r Redis) SetValueUntilChannelClose(key string, data string, ttl int, isOpen *bool) {
	go func() {
		/*
			We create a ticker for ttl - some small amount of time that we expect the delay will be when setting.
			This is so that as long as the app can set the key value pair, they will continue to exist in redis.
		*/
		ticker := time.NewTicker(time.Duration(ttl-2) * time.Second)
		defer func() {
			ticker.Stop()
			r.Delete(key)
		}()

		/*
			When the supplied boolean is not false yet, we set the value.
			Then we wait for the ticker to send a value before setting again
		*/
		for *isOpen {
			op1 := r.rdb.Set(context.Background(), key, data, time.Duration(ttl)*time.Second)
			r.lgr.NewInfo(fmt.Sprintf("[Redis][SetValueUntilChannelClose] On key %s this value was set: %s", key, data))
			if op1.Err() != nil {
				//Error
				log.Printf("Error when value %s is set", data)
				continue
			}
			<-ticker.C
		}
	}()
}

// GetValue attempts to obtain a value from redis with the specified key or return an error
func (r Redis) GetValue(key string) (res string, err error) {
	op2 := r.rdb.Get(context.Background(), key)
	if err = op2.Err(); err != nil {
		return
	}
	res, err = op2.Result()
	return
}
