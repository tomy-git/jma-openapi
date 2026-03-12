// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHTTPJMAClient_FetchWeatherAreaForecastDocument_AllowsMissingTemperatureSeries(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/forecast/data/forecast/130000.json" {
			http.NotFound(w, r)
			return
		}

		payload, err := loadForecastFixtureWithoutTemperatureSeries()
		if err != nil {
			t.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	client := NewHTTPJMAClient(server.URL)

	report, err := client.FetchWeatherAreaForecastDocument(context.Background(), "130000")
	if err != nil {
		t.Fatal(err)
	}

	if len(report.WeatherSeries) == 0 {
		t.Fatal("expected weather series")
	}
	if len(report.TempSeries) != 0 {
		t.Fatal("did not expect temperature series")
	}
}

func TestHTTPJMAClient_FetchForecastDocument_RequiresTemperatureSeries(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/forecast/data/forecast/130000.json" {
			http.NotFound(w, r)
			return
		}

		payload, err := loadForecastFixtureWithoutTemperatureSeries()
		if err != nil {
			t.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	client := NewHTTPJMAClient(server.URL)

	if _, err := client.FetchForecastDocument(context.Background(), "130000"); err == nil {
		t.Fatal("expected schema mismatch when temperature series is missing")
	}
}

func loadForecastFixtureWithoutTemperatureSeries() ([]byte, error) {
	payload, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "forecast-130000.json"))
	if err != nil {
		return nil, err
	}

	var reports []map[string]any
	if err := json.Unmarshal(payload, &reports); err != nil {
		return nil, err
	}

	timeSeries, ok := reports[0]["timeSeries"].([]any)
	if !ok {
		return nil, errors.New("timeSeries was not an array")
	}
	reports[0]["timeSeries"] = timeSeries[:2]

	return json.Marshal(reports)
}
