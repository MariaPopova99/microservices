package urls

import (
	"context"
	"errors"
	"time"

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

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.LongShortRepository {
	return &repo{db: db}
}

func (r *repo) ShortLongIns(ctx context.Context, inS *model.ShortUrls, inL *model.LongUrls) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(shortColumn, longColumn, createdAtColumn).
		Values(inS.ShortUrl, inL.LongUrl, time.Now().UTC()).
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

func (r *repo) GetLong(ctx context.Context, inS string) (*model.UrlFullInfo, error) {
	builder := sq.Select(idColumn, longColumn, shortColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{shortColumn: inS}).
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
			// Тут сделаем вызов ShortLongIns
			return nil, model.ErrorNoteNotFound
		}

		return nil, err
	}

	return &url, nil
}

func (r *repo) GetShort(ctx context.Context, inL string) (*model.UrlFullInfo, error) {
	builder := sq.Select(idColumn, longColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{longColumn: inL}).
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
			// Тут сделаем вызов ShortLongIns
			return nil, model.ErrorNoteNotFound
		}

		return nil, err
	}

	return &url, nil
}
