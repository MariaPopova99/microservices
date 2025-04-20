package urls

import (
	"context"

	"github.com/MariaPopova99/microservices/internal/converter"
	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) GetShort(ctx context.Context, req *desc.GetShortRequest) (*desc.GetShortResponse, error) {
	short_url, err := i.urlsService.GetShort(ctx, converter.ToLongUrlsFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &desc.GetShortResponse{
		Short_Url: short_url.ShortUrl,
		CreatedAt: timestamppb.New(short_url.CreatedAt),
	}, nil
}
