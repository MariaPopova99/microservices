package converter

import (
	"github.com/MariaPopova99/microservices/internal/model"
	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
)

func ToLongUrlsFromService(longUrl *model.LongUrls) *desc.GetShortRequest {

	return &desc.GetShortRequest{
		Long_Url: longUrl.LongUrl,
	}
}

func ToShortUrlsFromService(shortUrls model.ShortUrls) *desc.GetLongRequest {
	return &desc.GetLongRequest{
		Short_Url: shortUrls.ShortUrl,
	}
}

func ToLongUrlsFromDesc(longUrl *desc.GetShortRequest) *model.LongUrls {
	return &model.LongUrls{
		LongUrl: longUrl.Long_Url,
	}
}

func ToShortUrlsFromDesc(shortUrls *desc.GetLongRequest) *model.ShortUrls {
	return &model.ShortUrls{
		ShortUrl: shortUrls.Short_Url,
	}
}
