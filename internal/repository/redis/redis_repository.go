package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

// Save сохраняет оригинальный URL с коротким кодом как ключом
func (r *RedisRepository) Save(ctx context.Context, shortCode string, originalURL string) error {
	err := r.client.Set(ctx, shortCode, originalURL, 24*time.Hour).Err()
	if err != nil {
		log.Error().
			Err(err).
			Str("short_code", shortCode).
			Str("original_url", originalURL).
			Msg("failed to save URL to Redis")
		return err
	}

	log.Debug().
		Str("short_code", shortCode).
		Str("original_url", originalURL).
		Msg("URL saved to Redis")
	return nil
}

// Get получает оригинальный URL по shortCode
func (r *RedisRepository) Get(ctx context.Context, shortCode string) (string, error) {
	originalURL, err := r.client.Get(ctx, shortCode).Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("short_code", shortCode).
			Msg("URL not found in Redis")
		return "", err
	}

	log.Debug().
		Str("short_code", shortCode).
		Str("original_url", originalURL).
		Msg("URL retrieved from Redis")
	return originalURL, nil
}

// Delete удаляет запись по shortCode
func (r *RedisRepository) Delete(ctx context.Context, shortCode string) error {
	count, err := r.client.Del(ctx, shortCode).Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("short_code", shortCode).
			Msg("failed to delete from Redis")
		return err
	}
	if count == 0 {
		log.Warn().
			Str("short_code", shortCode).
			Msg("no URL found to delete in Redis")
		return redis.Nil
	}

	log.Debug().
		Str("short_code", shortCode).
		Msg("URL deleted from Redis")
	return nil
}
