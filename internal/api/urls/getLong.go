package urls

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/converter"
	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) GetLong(ctx context.Context, req *desc.GetLongRequest) (*desc.GetLongResponse, error) {
	long_url, err := i.urlsService.GetLong(ctx, converter.ToShortUrlsFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &desc.GetLongResponse{
		Long_Url:  long_url.ShortUrl,
		CreatedAt: timestamppb.New(long_url.CreatedAt),
	}, nil
}
