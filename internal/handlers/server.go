package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tomy-git/jma-openapi/internal/gen"
	"github.com/tomy-git/jma-openapi/internal/shared"
	"github.com/tomy-git/jma-openapi/internal/usecases"
)

type Server struct {
	healthUsecase   usecases.HealthUsecase
	areaUsecase     usecases.AreaUsecase
	forecastUsecase usecases.ForecastUsecase
}

func NewServer(
	healthUsecase usecases.HealthUsecase,
	areaUsecase usecases.AreaUsecase,
	forecastUsecase usecases.ForecastUsecase,
) *Server {
	return &Server{
		healthUsecase:   healthUsecase,
		areaUsecase:     areaUsecase,
		forecastUsecase: forecastUsecase,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeHandlerError(w http.ResponseWriter, r *http.Request, err error) {
	shared.WriteError(w, r, err)
}

var _ gen.ServerInterface = (*Server)(nil)
