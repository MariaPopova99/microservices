package repository

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/model"
)

type LongShortRepository interface {
	GetShort(ctx context.Context, LongUrl *model.LongUrls) (*model.UrlFullInfo, error)
	GetLong(ctx context.Context, shortUrl *model.ShortUrls) (*model.UrlFullInfo, error)
}
