package service

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/model"
)

type LongShortService interface {
	GetShort(ctx context.Context, longUrl *model.LongUrls) (*model.UrlFullInfo, error)
	GetLong(ctx context.Context, shortUrl *model.ShortUrls) (*model.UrlFullInfo, error)
}
