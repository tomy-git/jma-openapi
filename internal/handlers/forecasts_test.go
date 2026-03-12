// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/tomy-git/jma-openapi/internal/gen"
)

func TestGetV1ForecastsOfficeCodeAreasAreaCode(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/forecasts/130000/areas/130010", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if bytes.Contains(rec.Body.Bytes(), []byte("temperatureAreas")) {
		t.Fatal("did not expect temperatureAreas in forecast area response")
	}

	var response gen.ForecastAreaResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Office.Code != "130000" {
		t.Fatalf("expected office code 130000, got %s", response.Office.Code)
	}
	if response.WeatherArea.Code != "130010" {
		t.Fatalf("expected area code 130010, got %s", response.WeatherArea.Code)
	}
}

func TestGetV1ForecastsOfficeCodeAreasAreaCode_NotFound(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/forecasts/130000/areas/999999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	var response gen.ErrorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error.Code != "WEATHER_AREA_NOT_FOUND" {
		t.Fatalf("expected WEATHER_AREA_NOT_FOUND, got %s", response.Error.Code)
	}
}
