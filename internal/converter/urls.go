package converter

import (
	"github.com/MariaPopova99/microservices/internal/model"
	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
)

func ToLongUrlsFromService(longUrl *model.LongUrls) *desc.GetShortRequest {

	return &desc.GetShortRequest{
		LongUrl: longUrl.LongUrl,
	}
}

func ToShortUrlsFromService(shortUrls model.ShortUrls) *desc.GetLongRequest {
	return &desc.GetLongRequest{
		ShortUrl: shortUrls.ShortUrl,
	}
}

func ToLongUrlsFromDesc(longUrl *desc.GetShortRequest) *model.LongUrls {
	return &model.LongUrls{
		LongUrl: longUrl.LongUrl,
	}
}

func ToShortUrlsFromDesc(shortUrls *desc.GetLongRequest) *model.ShortUrls {
	return &model.ShortUrls{
		ShortUrl: shortUrls.ShortUrl,
	}
}
