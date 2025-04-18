package urls

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"crypto/md5"
	"encoding/hex"

	db "github.com/MariaPopova99/microservices/internal/db"
	"github.com/MariaPopova99/microservices/internal/model"
	repository "github.com/MariaPopova99/microservices/internal/repository/urls"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tableName = "urls"

	idColumn        = "id"
	shortColumn     = "short"
	longColumn      = "long"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

const (
	shortURLLength     = 8
	longURLExtraLength = 6
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.LongShortRepository {
	return &repo{db: db}
}

func (r *repo) CreateNewURL(ctx context.Context, shortUrl *model.ShortUrls, longUrl *model.LongUrls) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(shortColumn, longColumn, createdAtColumn).
		Values(shortUrl.ShortUrl, longUrl.LongUrl, time.Now().UTC()).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "note_repository.ShortLongIns",
		QueryRaw: query,
	}

	var id int64
	err = r.db.QueryRow(ctx, q.QueryRaw, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) GetLong(ctx context.Context, shortUrl *model.ShortUrls) (*model.UrlFullInfo, error) {
	builder := sq.Select(idColumn, longColumn, shortColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{shortColumn: shortUrl.ShortUrl}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "url_repository.GetLongUrl",
		QueryRaw: query,
	}

	var url model.UrlFullInfo
	err = r.db.QueryRow(ctx, q.QueryRaw, args...).Scan(&url.ID, &url.ShortUrl, &url.LongUrl, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			longUrl, err := generateLongUrl(shortUrl)
			if err != nil {
				return nil, err
			}
			id, err := r.CreateNewURL(ctx, shortUrl, longUrl)
			if err != nil {
				return nil, err
			}
			url.ID, url.ShortUrl, url.LongUrl, url.CreatedAt = id, shortUrl.ShortUrl, longUrl.LongUrl, time.Now()
		} else {
			return nil, err
		}
	}

	return &url, nil
}

func (r *repo) GetShort(ctx context.Context, longUrl *model.LongUrls) (*model.UrlFullInfo, error) {
	builder := sq.Select(idColumn, shortColumn, longColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{longColumn: longUrl.LongUrl}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "url_repository.GetShortUrl",
		QueryRaw: query,
	}

	var url model.UrlFullInfo
	err = r.db.QueryRow(ctx, q.QueryRaw, args...).Scan(&url.ID, &url.ShortUrl, &url.LongUrl, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			shortUrl, err := generateShortUrl(longUrl)
			if err != nil {
				return nil, err
			}
			id, err := r.CreateNewURL(ctx, shortUrl, longUrl)
			if err != nil {
				return nil, err
			}
			url.ID, url.ShortUrl, url.LongUrl, url.CreatedAt = id, shortUrl.ShortUrl, longUrl.LongUrl, time.Now()
		} else {
			return nil, err
		}

	}

	return &url, nil
}
func generateShortUrl(longUrl *model.LongUrls) (*model.ShortUrls, error) {
	// Создаём хеш от длинного URL
	hash := md5.Sum([]byte(longUrl.LongUrl))

	shortHash := hex.EncodeToString(hash[:])[:shortURLLength]
	return &model.ShortUrls{shortHash}, nil
}

func generateLongUrl(shortUrl *model.ShortUrls) (*model.LongUrls, error) {
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
