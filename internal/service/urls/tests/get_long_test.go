package tests

import (
	"context"
	"database/sql"
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

func TestGetLong(t *testing.T) {
	t.Parallel()
	type longShortRepositoryMockFunc func(mc *minimock.Controller) repository.LongShortRepository

	type args struct {
		ctx context.Context
		req *model.ShortUrls
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		shortUrl  = gofakeit.Animal()
		longUrl   = gofakeit.Animal()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		shortUrlModel = &model.ShortUrls{
			ShortUrl: shortUrl,
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
				req: shortUrlModel,
			},
			want: req,
			err:  nil,
			longShortRepositoryMockFunc: func(mc *minimock.Controller) repository.LongShortRepository {
				mock := repoMocks.NewLongShortRepositoryMock(mc)
				mock.GetLongMock.Expect(ctx, shortUrlModel).Return(req, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx: ctx,
				req: shortUrlModel,
			},
			want: nil,
			err:  repoErr,
			longShortRepositoryMockFunc: func(mc *minimock.Controller) repository.LongShortRepository {
				mock := repoMocks.NewLongShortRepositoryMock(mc)
				mock.GetLongMock.Expect(ctx, shortUrlModel).Return(nil, repoErr)
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

			resHandler, err := service.GetLong(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)

		})
	}

}
