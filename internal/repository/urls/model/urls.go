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
	ID        int64        `db:"id"`
	ShortUrl  string       `db:"short"`
	LongUrl   string       `db:"long"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
