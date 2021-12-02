package cache

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

type ClusterConfig struct {
	Addr     string // e.g. 10.0.0.1:6379,10.0.0.2:6379
	Password string
}

type Cluster struct {
	ClusterClient *redis.ClusterClient
}

func NewClusterClient(addr []string, password string) Cluster {
	return Cluster{
		ClusterClient: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addr,
			Password: password,
		}),
	}
}

func NewClusterWithCfg(cfg ClusterConfig) Cluster {
	return NewClusterClient(
		strings.Split(cfg.Addr, ","),
		cfg.Password,
	)
}

func (r *Cluster) Ping() error {
	if _, err := r.ClusterClient.Ping().Result(); err != nil {
		return err
	}
	return nil
}

func (r *Cluster) Set(key string, v interface{}, timeout time.Duration) error {
	j, err := json.Marshal(&v)
	if err != nil {
		return err
	}
	if _, err = r.ClusterClient.Set(key, j, timeout).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Cluster) Get(key string, v interface{}) error {
	b, err := r.ClusterClient.Get(key).Bytes()
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

func (r *Cluster) FindByPattern(pattern string) (map[string]interface{}, error) {
	var values map[string]interface{}
	keys, err := r.ClusterClient.Keys(pattern).Result()
	if err != nil {
		return values, err
	}
	for _, key := range keys {
		b, err := r.ClusterClient.Get(key).Bytes()
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, values[key]); err != nil {
			continue
		}
	}
	return values, nil
}

func (r *Cluster) Del(key string) error {
	return r.ClusterClient.Del(key).Err()
}

func (r *Cluster) DelByPattern(pattern string) error {
	keys, err := r.ClusterClient.Keys(pattern).Result()
	if err != nil {
		return err
	}
	return r.ClusterClient.Del(keys...).Err()
}

func (r *Cluster) Flush() error {
	return r.ClusterClient.FlushAll().Err()
}

func (r *Cluster) FlushByPattern(pattern string) error {
	return r.DelByPattern(pattern)
}

func (r *Cluster) IsKeyNotFound(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}

func (r *Cluster) SetExpire(key string, expiration time.Duration) bool {
	return r.ClusterClient.Expire(key, expiration).Val()
}

func (r *Cluster) SetExpireAt(key string, expireAt time.Time) bool {
	return r.ClusterClient.ExpireAt(key, expireAt).Val()
}
