package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisAuthorizationCodeRepository struct {
	rdb     *redis.Client
	ttlCode time.Duration
}

func NewRedisAuthorizationCodeRepository(rdb *redis.Client, ttlCode time.Duration) *RedisAuthorizationCodeRepository {
	return &RedisAuthorizationCodeRepository{rdb, ttlCode}
}

func (r *RedisAuthorizationCodeRepository) CheckEndDelete(
	ctx context.Context,
	code string,
) (bool, error) {
	result, err := r.rdb.GetDel(ctx, code).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if result != "" {
		return true, nil
	}

	return false, nil
}

func (r *RedisAuthorizationCodeRepository) Save(
	ctx context.Context,
	code string,
) error {
	status := r.rdb.Set(ctx, code, true, r.ttlCode)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}
