package repository

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/model"
)

type LongShortRepository interface {
	GetShort(ctx context.Context, inS string) (*model.UrlFullInfo, error)
	GetLong(ctx context.Context, inL string) (*model.UrlFullInfo, error)
}
