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
	if response.WeatherArea == nil {
		t.Fatal("expected weather area")
	}
	if response.WeatherArea.Code != "130010" {
		t.Fatalf("expected area code 130010, got %s", response.WeatherArea.Code)
	}
	if response.TemperatureArea != nil {
		t.Fatal("did not expect temperature area")
	}
}

func TestGetV1ForecastsOfficeCodeAreasAreaCode_TemperatureArea(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/forecasts/130000/areas/44132", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if bytes.Contains(rec.Body.Bytes(), []byte("weatherArea")) {
		t.Fatal("did not expect weatherArea in temperature area response")
	}

	var response gen.ForecastAreaResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.TemperatureArea == nil {
		t.Fatal("expected temperature area")
	}
	if response.TemperatureArea.Code != "44132" {
		t.Fatalf("expected area code 44132, got %s", response.TemperatureArea.Code)
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

	if response.Error.Code != "FORECAST_AREA_NOT_FOUND" {
		t.Fatalf("expected FORECAST_AREA_NOT_FOUND, got %s", response.Error.Code)
	}
}

func TestGetV1ForecastsOfficeCodeAreas(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/forecasts/130000/areas", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response gen.ForecastAreaListResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if len(response.Items) != 8 {
		t.Fatalf("expected 8 items, got %d", len(response.Items))
	}
}

func TestGetV1ForecastsOfficeCodeWeatherAreas(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/forecasts/130000/weather-areas", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response gen.WeatherAreaListResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if len(response.Items) != 4 || response.Items[0].Kind != gen.Weather {
		t.Fatalf("expected weather refs, got %+v", response.Items)
	}
}

func TestGetV1ForecastsOfficeCodeAreasResolve(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/forecasts/130000/areas:resolve?q=%E6%9D%B1%E4%BA%AC&matchMode=suggested", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response gen.ForecastAreaListResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if len(response.Items) == 0 {
		t.Fatal("expected resolve candidates")
	}
}

func TestGetV1ForecastsOfficeCodeAreasAreaCodeLatest(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/forecasts/130000/areas/44132/latest", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response gen.ForecastAreaResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.TemperatureArea == nil || len(response.TemperatureArea.TimeSeries) != 1 {
		t.Fatalf("expected single latest temperature entry, got %+v", response.TemperatureArea)
	}
}
