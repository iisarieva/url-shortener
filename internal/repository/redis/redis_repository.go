package redis

import (
	"context"
	"errors"
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
			Msg("❌ Redis: не удалось сохранить ссылку")
		return err
	}

	log.Debug().
		Str("short_code", shortCode).
		Str("original_url", originalURL).
		Msg("💾 Redis: ссылка сохранена")
	return nil
}

// Get получает оригинальный URL по shortCode
func (r *RedisRepository) Get(ctx context.Context, shortCode string) (string, error) {
	originalURL, err := r.client.Get(ctx, shortCode).Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("short_code", shortCode).
			Msg("❌ Redis: не удалось получить ссылку")
		return "", err
	}

	log.Debug().
		Str("short_code", shortCode).
		Str("original_url", originalURL).
		Msg("📦 Redis: ссылка получена")
	return originalURL, nil
}

// Delete удаляет запись по shortCode
func (r *RedisRepository) Delete(ctx context.Context, shortCode string) error {
	count, err := r.client.Del(ctx, shortCode).Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("short_code", shortCode).
			Msg("❌ Redis: ошибка при удалении")
		return err
	}
	if count == 0 {
		log.Warn().
			Str("short_code", shortCode).
			Msg("⚠️ Redis: ссылка не найдена для удаления")
		return errors.New("not found")
	}

	log.Debug().
		Str("short_code", shortCode).
		Msg("🗑️ Redis: ссылка удалена")
	return nil
}
