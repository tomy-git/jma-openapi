// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
	"github.com/tomy-git/jma-openapi/internal/handlers"
	"github.com/tomy-git/jma-openapi/internal/mappers"
	"github.com/tomy-git/jma-openapi/internal/shared"
	"github.com/tomy-git/jma-openapi/internal/usecases"
)

func main() {
	cfg := shared.LoadConfig()
	logger := shared.NewLogger(cfg.LogFormat)

	client := clients.NewHTTPJMAClient(cfg.JMABaseURL)
	areaMapper := mappers.NewAreaMapper()
	forecastMapper := mappers.NewForecastMapper()

	server := handlers.NewServer(
		usecases.NewHealthUsecase(cfg.Service, cfg.Version),
		usecases.NewAreaUsecase(client, areaMapper),
		usecases.NewForecastUsecase(client, areaMapper, forecastMapper),
	)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(shared.RequestLoggingMiddleware(logger))

	gen.RegisterHandlers(router, server)
	registerOpenAPIRoutes(router)

	address := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("server starting", "address", address)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal(err)
	}
}

func registerOpenAPIRoutes(router chi.Router) {
	router.Get("/assets/scalar-api-reference.js", func(w http.ResponseWriter, r *http.Request) {
		payload, err := os.ReadFile("web/scalar-api-reference.js")
		if err != nil {
			shared.WriteError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write(payload)
	})

	router.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		payload, err := gen.LoadOpenAPISpecYAML()
		if err != nil {
			shared.WriteError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(payload)
	})

	router.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		payload, err := gen.LoadOpenAPISpecJSON()
		if err != nil {
			shared.WriteError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	})

	router.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html>
<html lang="ja">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>jma-openapi docs</title>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.yaml"></script>
    <script src="/assets/scalar-api-reference.js"></script>
  </body>
</html>`))
	})
}
