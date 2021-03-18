package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Password string
	Db       int
}

type Redis struct {
	Client *redis.Client
}

func New(host string, port string, password string, db int) Redis {
	return Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       db,
		}),
	}
}

func NewWithCfg(cfg Config) Redis {
	return New(
		cfg.Host,
		cfg.Port,
		cfg.Password,
		cfg.Db,
	)
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

func (r *Redis) FindByPattern(pattern string) (map[string]interface{}, error) {
	var values map[string]interface{}
	keys, err := r.Client.Keys(pattern).Result()
	if err != nil {
		return values, err
	}
	for _, key := range keys {
		b, err := r.Client.Get(key).Bytes()
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, values[key]); err != nil {
			continue
		}
	}
	return values, nil
}

func (r *Redis) Del(key string) error {
	return r.Client.Del(key).Err()
}

func (r *Redis) DelByPattern(pattern string) error {
	keys, err := r.Client.Keys(pattern).Result()
	if err != nil {
		return err
	}
	return r.Client.Del(keys...).Err()
}

func (r *Redis) Flush() error {
	return r.Client.FlushAll().Err()
}

func (r *Redis) FlushByPattern(pattern string) error {
	return r.DelByPattern(pattern)
}

func (r *Redis) IsKeyNotFound(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}

func (r *Redis) SetExpire(key string, expiration time.Duration) bool {
	return r.Client.Expire(key, expiration).Val()
}

func (r *Redis) SetExpireAt(key string, expireAt time.Time) bool {
	return r.Client.ExpireAt(key, expireAt).Val()
}
