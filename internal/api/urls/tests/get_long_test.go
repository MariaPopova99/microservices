package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/MariaPopova99/microservices/internal/api/urls"
	"github.com/MariaPopova99/microservices/internal/model"
	"github.com/MariaPopova99/microservices/internal/service"
	serviceMocks "github.com/MariaPopova99/microservices/internal/service/mocks"
	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetLong(t *testing.T) {
	type longShortServiceMockFunc func(mc *minimock.Controller) service.LongShortService
	type args struct {
		ctx context.Context
		req *desc.GetLongRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		longUrl   = gofakeit.Animal()
		shortUrl  = gofakeit.Animal()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetLongRequest{
			ShortUrl: shortUrl,
		}

		res = &desc.GetLongResponse{
			LongUrl:   longUrl,
			CreatedAt: timestamppb.New(createdAt),
		}

		shortUrlModel = &model.ShortUrls{
			ShortUrl: shortUrl,
		}

		fullInfo = &model.UrlFullInfo{
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
		name                     string
		args                     args
		want                     *desc.GetLongResponse
		err                      error
		longShortServiceMockFunc longShortServiceMockFunc
	}{
		{
			name: "success test service",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			longShortServiceMockFunc: func(mc *minimock.Controller) service.LongShortService {
				mock := serviceMocks.NewLongShortServiceMock(mc)
				mock.GetLongMock.Expect(ctx, shortUrlModel).Return(fullInfo, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			longShortServiceMockFunc: func(mc *minimock.Controller) service.LongShortService {
				mock := serviceMocks.NewLongShortServiceMock(mc)
				mock.GetLongMock.Expect(ctx, shortUrlModel).Return(nil, serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			longShortServiceMockFunc := tt.longShortServiceMockFunc(mc)
			api := urls.NewImplementation(longShortServiceMockFunc)

			resHandler, err := api.GetLong(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)

		})
	}
}
