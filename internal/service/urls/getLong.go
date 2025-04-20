package urls

import (
	"context"
	"math/rand"

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

func GenerateLongUrl(shortUrl *model.ShortUrls) (*model.LongUrls, error) {
	randomString := func(length int) string {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		longUrl := make([]byte, length)
		for i := range longUrl {
			longUrl[i] = charset[rand.Intn(len(charset))]
		}
		return string(longUrl)
	}
	longUrl := shortUrl.ShortUrl + randomString(longURLExtraLength)
	return &model.LongUrls{longUrl}, nil
}
