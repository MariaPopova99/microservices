package urls

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/MariaPopova99/microservices/internal/model"
)

const (
	shortURLLength = 8
)

func (s *serv) GetShort(ctx context.Context, longUrl *model.LongUrls) (*model.UrlFullInfo, error) {
	log.Printf("GetShort server started: %s", longUrl.LongUrl)
	shortUrl, err := s.urlRepository.GetShort(ctx, longUrl)
	log.Printf("GetShort server started err: %s", err)
	if err != nil {
		if errors.Is(err, model.ErrorNoteNotFound) {
			log.Printf("GetShort server started err dd: %s", err)
			shortUrl, err := GenerateShortUrl(longUrl)
			if err != nil {
				return nil, err
			}
			id, err := s.urlRepository.CreateNewURL(ctx, shortUrl, longUrl)
			if err != nil {
				return nil, err
			}
			return &model.UrlFullInfo{
				ID:        id,
				ShortUrl:  shortUrl.ShortUrl,
				LongUrl:   longUrl.LongUrl,
				CreatedAt: time.Now(),
				UpdatedAt: sql.NullTime{Time: time.Now(), Valid: false},
			}, nil
		}
		return nil, err
	}
	return shortUrl, nil
}

func GenerateShortUrl(longUrl *model.LongUrls) (*model.ShortUrls, error) {
	// Создаём хеш от длинного URL
	hash := md5.Sum([]byte(longUrl.LongUrl))

	shortHash := hex.EncodeToString(hash[:])[:shortURLLength]
	return &model.ShortUrls{ShortUrl: shortHash}, nil
}
