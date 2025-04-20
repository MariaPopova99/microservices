package urls

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/model"
)

func (s *serv) GetShort(ctx context.Context, longUrl *model.LongUrls) (*model.UrlFullInfo, error) {
	shortURL, err := s.urlRepository.GetShort(ctx, longUrl)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}
