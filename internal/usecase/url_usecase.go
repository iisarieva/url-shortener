package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"url-shortener/internal/domain/repository"

	"github.com/rs/zerolog/log"
)

type URLUseCase struct {
	repo repository.URLRepository
}

func NewURLUseCase(r repository.URLRepository) *URLUseCase {
	return &URLUseCase{repo: r}
}

// CreateShortURL генерирует короткий код, сохраняет его в Redis
func (u *URLUseCase) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	start := time.Now()

	shortCode := generateShortCode()

	err := u.repo.Save(ctx, shortCode, originalURL)
	if err != nil {
		log.Error().
			Err(err).
			Str("short_code", shortCode).
			Str("original_url", originalURL).
			Msg("❌ Не удалось сохранить короткую ссылку")
		return "", err
	}
	log.Info().
		Str("short_code", shortCode).
		Str("original_url", originalURL).
		Dur("duration", time.Since(start)).
		Msg("✅ Сокращена ссылка")
	return shortCode, nil
}

// GetOriginalURL достаёт оригинальный URL по shortCode
func (u *URLUseCase) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	originalURL, err := u.repo.Get(ctx, shortCode)
	if err != nil {
		log.Error().
			Err(err).
			Str("short_code", shortCode).
			Msg("❌ Ошибка при получении оригинального URL")
		return "", err
	}

	log.Info().
		Str("short_code", shortCode).
		Str("original_url", originalURL).
		Msg("🔗 Получен оригинальный URL")

	return originalURL, nil
}

// DeleteShortURL удаляет ссылку по shortCode
func (u *URLUseCase) DeleteShortURL(ctx context.Context, shortCode string) error {
	err := u.repo.Delete(ctx, shortCode)
	if err != nil {
		log.Warn().
			Err(err).
			Str("short_code", shortCode).
			Msg("⚠️ Ссылка не найдена или уже удалена")
		return err
	}

	log.Info().
		Str("short_code", shortCode).
		Msg("🗑️ Ссылка удалена")

	return nil
}

// Вспомогательная функция: генерирует случайный short code
func generateShortCode() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
