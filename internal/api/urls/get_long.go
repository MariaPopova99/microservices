package urls

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/converter"
	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) GetLong(ctx context.Context, req *desc.GetLongRequest) (*desc.GetLongResponse, error) {
	longUrl, err := i.urlsService.GetLong(ctx, converter.ToShortUrlsFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &desc.GetLongResponse{
		LongUrl:   longUrl.LongUrl,
		CreatedAt: timestamppb.New(longUrl.CreatedAt),
	}, nil
}
