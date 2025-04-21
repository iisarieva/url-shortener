package repository

import "context"

type URLRepository interface {
	Save(ctx context.Context, shortCode string, originalURL string) error
	Get(ctx context.Context, shortCode string) (string, error)
	Delete(ctx context.Context, shortCode string) error
}
