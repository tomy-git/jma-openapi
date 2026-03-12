// Code generated for the current OpenAPI contract scaffold. DO NOT EDIT.

package gen

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetV1AreasParams struct {
	Parent *string
}

type ServerInterface interface {
	GetHealthz(w http.ResponseWriter, r *http.Request)
	GetV1Areas(w http.ResponseWriter, r *http.Request, params GetV1AreasParams)
	GetV1AreasAreaCode(w http.ResponseWriter, r *http.Request, areaCode string)
	GetV1ForecastsOfficeCode(w http.ResponseWriter, r *http.Request, officeCode string)
}

func RegisterHandlers(router chi.Router, server ServerInterface) {
	router.Get("/healthz", server.GetHealthz)
	router.Get("/v1/areas", func(w http.ResponseWriter, r *http.Request) {
		var params GetV1AreasParams
		if parent := r.URL.Query().Get("parent"); parent != "" {
			params.Parent = &parent
		}

		server.GetV1Areas(w, r, params)
	})
	router.Get("/v1/areas/{areaCode}", func(w http.ResponseWriter, r *http.Request) {
		server.GetV1AreasAreaCode(w, r, chi.URLParam(r, "areaCode"))
	})
	router.Get("/v1/forecasts/{officeCode}", func(w http.ResponseWriter, r *http.Request) {
		server.GetV1ForecastsOfficeCode(w, r, chi.URLParam(r, "officeCode"))
	})
}
