package urls

import (
	repository "github.com/MariaPopova99/microservices/internal/repository"
	service "github.com/MariaPopova99/microservices/internal/service"
)

type serv struct {
	urlRepository repository.LongShortRepository
}

func NewService(urlRepository repository.LongShortRepository) service.LongShortService {
	return &serv{
		urlRepository: urlRepository,
	}
}
