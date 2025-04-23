package urls

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/MariaPopova99/microservices/internal/db"
	"github.com/MariaPopova99/microservices/internal/model"
	"github.com/MariaPopova99/microservices/internal/repository"
	"github.com/MariaPopova99/microservices/internal/repository/urls/converter"
	repoModel "github.com/MariaPopova99/microservices/internal/repository/urls/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName = "urls"

	idColumn        = "id"
	shortColumn     = "short"
	longColumn      = "long"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
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
	log.Printf("GetLong starting: %s", shortUrl.ShortUrl)
	builder := sq.Select(idColumn, longColumn, shortColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{shortColumn: shortUrl.ShortUrl}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("builder failed: %s", err)
		return nil, err
	}

	q := db.Query{
		Name:     "url_repository.GetLongUrl",
		QueryRaw: query,
	}

	var url repoModel.UrlFullInfo
	err = pgxscan.Get(ctx, r.db, &url, q.QueryRaw, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUrlFullInfoFromRepo(&url), nil
}

func (r *repo) GetShort(ctx context.Context, longUrl *model.LongUrls) (*model.UrlFullInfo, error) {
	log.Printf("GetShort starting: %s", longUrl.LongUrl)
	builder := sq.Select(idColumn, shortColumn, longColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{longColumn: longUrl.LongUrl}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("builder failed: %s", err)
		return nil, err
	}

	q := db.Query{
		Name:     "url_repository.GetShortUrl",
		QueryRaw: query,
	}

	var url repoModel.UrlFullInfo
	err = pgxscan.Get(ctx, r.db, &url, q.QueryRaw, args...)

	if err != nil {
		if isNoRowsError(err) {
			return nil, model.ErrorNoteNotFound
		}
		log.Printf("QueryRow failed: %s", err)

		return nil, err

	}

	return converter.ToUrlFullInfoFromRepo(&url), nil
}

func isNoRowsError(err error) bool {
	// Прямое сравнение с pgx.ErrNoRows
	if errors.Is(err, pgx.ErrNoRows) {
		return true
	}

	// Проверка через PgError
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "02000" // Код "no data found"
	}

	// Проверка текста
	if strings.Contains(strings.ToLower(err.Error()), "no rows") {
		return true
	}

	return false
}
