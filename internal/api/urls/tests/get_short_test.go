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

func TestGetShort(t *testing.T) {
	type longShortServiceMockFunc func(mc *minimock.Controller) service.LongShortService
	type args struct {
		ctx context.Context
		req *desc.GetShortRequest
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

		req = &desc.GetShortRequest{
			LongUrl: longUrl,
		}

		longUrlModel = &model.LongUrls{
			LongUrl: longUrl,
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

		res = &desc.GetShortResponse{
			ShortUrl:  shortUrl,
			CreatedAt: timestamppb.New(createdAt),
		}
	)

	tests := []struct {
		name                     string
		args                     args
		want                     *desc.GetShortResponse
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
				mock.GetShortMock.Expect(ctx, longUrlModel).Return(fullInfo, nil) // Описываем как себя должен вести метод сервисного слоя.
				return mock                                                       // Вход и выход модельки для сервиса
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
				mock.GetShortMock.Expect(ctx, longUrlModel).Return(nil, serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			longShortServiceMockFunc := tt.longShortServiceMockFunc(mc)
			api := urls.NewImplementation(longShortServiceMockFunc) //Регистрируем апи слой

			resHadler, err := api.GetShort(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHadler) //Ждем от апи слоя протобафную структуру
		})
	}
}
