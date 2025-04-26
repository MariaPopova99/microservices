package urls

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/model"
)

const (
	longURLExtraLength = 6
)

func (s *serv) GetLong(ctx context.Context, shortUrl *model.ShortUrls) (*model.UrlFullInfo, error) {
	longURL, err := s.urlRepository.GetLong(ctx, shortUrl)
	if err != nil {
		return nil, err
	}
	return longURL, nil
}
