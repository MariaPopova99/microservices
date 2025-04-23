package urls

import (
	"github.com/MariaPopova99/microservices/internal/service"
	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
)

type Implementation struct {
	desc.UnimplementedLongShortV1Server
	urlsService service.LongShortService
}

func NewImplementation(urlsSevice service.LongShortService) *Implementation {
	return &Implementation{
		urlsService: urlsSevice,
	}
}
