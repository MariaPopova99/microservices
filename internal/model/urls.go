package model

import (
	"database/sql"
	"time"
)

type LongUrls struct {
	LongUrl string
}

type ShortUrls struct {
	ShortUrl string
}

type UrlInfo struct {
	Url       string
	CreatedAt time.Time
}

type UrlFullInfo struct {
	ID        int64
	ShortUrl  string
	LongUrl   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
