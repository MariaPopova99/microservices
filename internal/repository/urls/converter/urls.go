package converter

import (
	repoModel "github.com/MariaPopova99/microservices/internal/repository/urls/model"

	"github.com/MariaPopova99/microservices/internal/model"
)

func ToLongUrlsFromRepo(url *repoModel.LongUrls) *model.LongUrls {
	return &model.LongUrls{
		LongUrl: url.LongUrl,
	}
}

func ToShortUrlsFromRepo(url *repoModel.ShortUrls) *model.ShortUrls {
	return &model.ShortUrls{
		ShortUrl: url.ShortUrl,
	}
}

func ToUrlInfoFromRepo(url *repoModel.UrlInfo) *model.UrlInfo {
	return &model.UrlInfo{
		Url:       url.Url,
		CreatedAt: url.CreatedAt,
	}
}

func ToUrlFullInfoFromRepo(url *repoModel.UrlFullInfo) *model.UrlFullInfo {
	return &model.UrlFullInfo{
		ID:        url.ID,
		ShortUrl:  url.ShortUrl,
		LongUrl:   url.LongUrl,
		CreatedAt: url.CreatedAt,
		UpdatedAt: url.UpdatedAt,
	}
}
