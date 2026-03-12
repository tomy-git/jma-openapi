// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
	"github.com/tomy-git/jma-openapi/internal/mappers"
	"github.com/tomy-git/jma-openapi/internal/usecases"
)

func TestGetV1Areas_QueryFilters(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/areas?name=%E6%9D%B1%E4%BA%AC%E9%83%BD&officeName=%E6%B0%97%E8%B1%A1%E5%BA%81&child=130010", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response gen.AreasResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(response.Items))
	}
	if response.Items[0].Code != "130000" {
		t.Fatalf("expected code 130000, got %s", response.Items[0].Code)
	}
}

func TestGetV1Areas_EmptyItemsForNoMatch(t *testing.T) {
	t.Parallel()

	server := newTestServer(t)
	router := chi.NewRouter()
	gen.HandlerFromMux(server, router)

	req := httptest.NewRequest(http.MethodGet, "/v1/areas?name=%E8%A9%B2%E5%BD%93%E3%81%AA%E3%81%97", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response gen.AreasResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Items) != 0 {
		t.Fatalf("expected empty items, got %d", len(response.Items))
	}
}

func newTestServer(t *testing.T) *Server {
	t.Helper()

	client := fixtureClient{
		areaDocument:     loadAreaDocument(t),
		forecastDocument: loadForecastDocument(t),
	}

	return NewServer(
		usecases.NewHealthUsecase("jma-openapi", "test"),
		usecases.NewAreaUsecase(client, mappers.NewAreaMapper()),
		usecases.NewForecastUsecase(client, mappers.NewAreaMapper(), mappers.NewForecastMapper()),
	)
}

type fixtureClient struct {
	areaDocument     clients.AreaJSONDocument
	forecastDocument clients.ForecastReportJSON
}

func (c fixtureClient) FetchAreaDocument(context.Context) (clients.AreaJSONDocument, error) {
	return c.areaDocument, nil
}

func (c fixtureClient) FetchForecastDocument(context.Context, string) (clients.ForecastReportJSON, error) {
	return c.forecastDocument, nil
}

func (c fixtureClient) FetchWeatherAreaForecastDocument(context.Context, string) (clients.ForecastReportJSON, error) {
	return c.forecastDocument, nil
}

func loadAreaDocument(t *testing.T) clients.AreaJSONDocument {
	t.Helper()

	payload, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "area.json"))
	if err != nil {
		t.Fatal(err)
	}

	var document clients.AreaJSONDocument
	if err := json.Unmarshal(payload, &document); err != nil {
		t.Fatal(err)
	}

	return document
}

func loadForecastDocument(t *testing.T) clients.ForecastReportJSON {
	t.Helper()

	payload, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "forecast-130000.json"))
	if err != nil {
		t.Fatal(err)
	}

	var reports []map[string]json.RawMessage
	if err := json.Unmarshal(payload, &reports); err != nil {
		t.Fatal(err)
	}

	var report clients.ForecastReportJSON
	if err := json.Unmarshal(reports[0]["publishingOffice"], &report.PublishingOffice); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(reports[0]["reportDatetime"], &report.ReportDatetime); err != nil {
		t.Fatal(err)
	}

	var series []json.RawMessage
	if err := json.Unmarshal(reports[0]["timeSeries"], &series); err != nil {
		t.Fatal(err)
	}

	var weather clients.ForecastTimeSeriesWeatherJSON
	var pops clients.ForecastTimeSeriesPopJSON
	var temps clients.ForecastTimeSeriesTempJSON
	if err := json.Unmarshal(series[0], &weather); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(series[1], &pops); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(series[2], &temps); err != nil {
		t.Fatal(err)
	}

	report.WeatherSeries = []clients.ForecastTimeSeriesWeatherJSON{weather}
	report.PopSeries = []clients.ForecastTimeSeriesPopJSON{pops}
	report.TempSeries = []clients.ForecastTimeSeriesTempJSON{temps}

	return report
}
