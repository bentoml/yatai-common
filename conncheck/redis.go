package conncheck

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Addr     string
	Password string
}

type RedisProbe struct {
	client *redis.Client
}

func NewRedisProbe(cfg RedisConfig) *RedisProbe {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})

	return &RedisProbe{client: client}
}

func (p *RedisProbe) Test(ctx context.Context, prefix string) error {
	testKey := fmt.Sprintf("%s:%d", prefix, time.Now().UnixNano())
	testValue := "test-value"

	logrus.Infof("Testing Redis connection to %s", p.client.Options().Addr)
	err := p.client.Set(ctx, testKey, testValue, 1*time.Minute).Err()
	if err != nil {
		logrus.Errorf("Redis write test failed: %v", err)
		return fmt.Errorf("redis write test failed: %w", err)
	}
	logrus.Infof("Redis write test passed")

	val, err := p.client.Get(ctx, testKey).Result()
	if err != nil {
		logrus.Errorf("Redis read test failed: %v", err)
		return fmt.Errorf("redis read test failed: %w", err)
	}
	logrus.Infof("Redis read test passed")

	if val != testValue {
		logrus.Errorf("Redis value mismatch: got %s, want %s", val, testValue)
		return fmt.Errorf("redis value mismatch: got %s, want %s", val, testValue)
	}

	err = p.client.Del(ctx, testKey).Err()
	if err != nil {
		logrus.Errorf("Redis cleanup failed: %v", err)
		return fmt.Errorf("redis cleanup failed: %w", err)
	}
	logrus.Infof("Redis cleanup test passed")

	return nil
}

func (p *RedisProbe) Close() error {
	return p.client.Close()
}
