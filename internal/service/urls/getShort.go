package urls

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/MariaPopova99/microservices/internal/model"
)

const (
	shortURLLength = 8
)

func (s *serv) GetShort(ctx context.Context, longUrl *model.LongUrls) (*model.UrlFullInfo, error) {
	shortURL, err := s.urlRepository.GetShort(ctx, longUrl)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}

func GenerateShortUrl(longUrl *model.LongUrls) (*model.ShortUrls, error) {
	// Создаём хеш от длинного URL
	hash := md5.Sum([]byte(longUrl.LongUrl))

	shortHash := hex.EncodeToString(hash[:])[:shortURLLength]
	return &model.ShortUrls{shortHash}, nil
}
