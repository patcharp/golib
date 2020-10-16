package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type Config struct {
	Host string
	Port string
	Db   int
}

type Redis struct {
	Client *redis.Client
}

func New(host string, port string, db int) Redis {
	return Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: "",
			DB:       db,
		}),
	}
}

func (r *Redis) Ping() error {
	if _, err := r.Client.Ping().Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Set(key string, v interface{}, timeout time.Duration) error {
	j, err := json.Marshal(&v)
	if err != nil {
		return err
	}
	if _, err = r.Client.Set(key, j, timeout).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key string, v interface{}) error {
	b, err := r.Client.Get(key).Bytes()
	if err != nil {
		return err
	}
	if v != nil {
		if err := json.Unmarshal(b, v); err != nil {
			return err
		}
	}
	return nil
}

func (r *Redis) Del(key string) error {
	return r.Client.Del(key).Err()
}

func (r *Redis) DelPattern(key string) error {
	keys, err := r.Client.Keys(key).Result()
	if err != nil {
		return err
	}
	return r.Client.Del(keys...).Err()
}

func (r *Redis) Flush() error {
	return r.Client.FlushAll().Err()
}

func (r *Redis) IsKeyNotFound(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}
