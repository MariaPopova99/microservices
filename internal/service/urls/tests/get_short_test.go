package tests

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/MariaPopova99/microservices/internal/model"
	"github.com/MariaPopova99/microservices/internal/repository"
	repoMocks "github.com/MariaPopova99/microservices/internal/repository/mocks"
	"github.com/MariaPopova99/microservices/internal/service/urls"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetShort(t *testing.T) {
	t.Parallel()
	type longShortRepositoryMockFunc func(mc *minimock.Controller) repository.LongShortRepository

	type args struct {
		ctx context.Context
		req *model.LongUrls
	}
	const (
		shortURLLength = 8
	)
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		shortUrl  = gofakeit.Animal()
		longUrl   = gofakeit.Animal()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		longUrlModel = &model.LongUrls{
			LongUrl: longUrl,
		}

		repoErr = fmt.Errorf("repo error")

		req = &model.UrlFullInfo{
			ID:        id,
			ShortUrl:  shortUrl,
			LongUrl:   longUrl,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}

		hash           = md5.Sum([]byte(longUrl))
		shortHash      = hex.EncodeToString(hash[:])[:shortURLLength]
		shortHashModel = &model.ShortUrls{ShortUrl: shortHash}

		reqNewUrl = &model.UrlFullInfo{
			ID:        id,
			ShortUrl:  shortHash,
			LongUrl:   longUrl,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: false,
			},
		}
	)

	tests := []struct {
		name                        string
		args                        args
		want                        *model.UrlFullInfo
		err                         error
		longShortRepositoryMockFunc longShortRepositoryMockFunc
	}{
		{
			name: "success case repo",
			args: args{
				ctx: ctx,
				req: longUrlModel,
			},
			want: req,
			err:  nil,
			longShortRepositoryMockFunc: func(mc *minimock.Controller) repository.LongShortRepository {
				mock := repoMocks.NewLongShortRepositoryMock(mc)
				mock.GetShortMock.Expect(ctx, longUrlModel).Return(req, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx: ctx,
				req: longUrlModel,
			},
			want: nil,
			err:  repoErr,
			longShortRepositoryMockFunc: func(mc *minimock.Controller) repository.LongShortRepository {
				mock := repoMocks.NewLongShortRepositoryMock(mc)
				mock.GetShortMock.Expect(ctx, longUrlModel).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "repo success case create new short url",
			args: args{
				ctx: ctx,
				req: longUrlModel,
			},
			want: reqNewUrl,
			err:  nil,
			longShortRepositoryMockFunc: func(mc *minimock.Controller) repository.LongShortRepository {
				mock := repoMocks.NewLongShortRepositoryMock(mc)
				mock.GetShortMock.Expect(ctx, longUrlModel).Return(nil, model.ErrorNoteNotFound)
				mock.CreateNewURLMock.Expect(ctx, shortHashModel, longUrlModel).Return(id, nil)
				return mock
			},
		},
		{
			name: "repo error case create new url",
			args: args{
				ctx: ctx,
				req: longUrlModel,
			},
			want: nil,
			err:  repoErr,
			longShortRepositoryMockFunc: func(mc *minimock.Controller) repository.LongShortRepository {
				mock := repoMocks.NewLongShortRepositoryMock(mc)
				mock.GetShortMock.Expect(ctx, longUrlModel).Return(nil, model.ErrorNoteNotFound)
				mock.CreateNewURLMock.Expect(ctx, shortHashModel, longUrlModel).Return(0, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			shortRepoMock := tt.longShortRepositoryMockFunc(mc)
			service := urls.NewService(shortRepoMock)

			resHandler, err := service.GetShort(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			if resHandler != nil {
				require.Equal(t, tt.want.ID, resHandler.ID)
				require.Equal(t, tt.want.LongUrl, resHandler.LongUrl)
				require.Equal(t, tt.want.ShortUrl, resHandler.ShortUrl)
				require.False(t, resHandler.CreatedAt.IsZero())
				require.Equal(t, tt.want.UpdatedAt.Valid, resHandler.UpdatedAt.Valid)
			} else {
				require.Equal(t, tt.want, resHandler)
			}

		})
	}

}
