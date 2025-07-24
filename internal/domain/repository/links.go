package repository

import (
	"context"

	"github.com/onemoreslacker/shortener/internal/domain/entity"
)

type Links interface {
	Get(ctx context.Context, shortURL string) (*entity.Link, error)
	Add(ctx context.Context, link *entity.Link) error
}
