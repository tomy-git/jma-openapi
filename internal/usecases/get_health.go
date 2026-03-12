package usecases

import (
	"time"

	"github.com/tomy-git/jma-openapi/internal/gen"
)

type HealthUsecase struct {
	service string
	version string
}

func NewHealthUsecase(service, version string) HealthUsecase {
	return HealthUsecase{
		service: service,
		version: version,
	}
}

func (u HealthUsecase) Execute() gen.HealthResponse {
	return gen.HealthResponse{
		Status:    "ok",
		Service:   u.service,
		Version:   u.version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
