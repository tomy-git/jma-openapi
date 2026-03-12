// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package mappers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
)

func TestForecastMapper_ToForecastResponse(t *testing.T) {
	t.Parallel()

	payload, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "forecast-130000.json"))
	if err != nil {
		t.Fatal(err)
	}

	var reports []map[string]json.RawMessage
	if err := json.Unmarshal(payload, &reports); err != nil {
		t.Fatal(err)
	}
	if len(reports) == 0 {
		t.Fatal("expected forecast reports")
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

	mapper := NewForecastMapper()
	response := mapper.ToForecastResponse(report, gen.Area{Code: "130000", Name: "東京都"})

	if response.Office.Code != "130000" {
		t.Fatalf("expected office code 130000, got %s", response.Office.Code)
	}
	if len(response.WeatherAreas) == 0 {
		t.Fatal("expected weather areas")
	}
	if len(response.TemperatureAreas) == 0 {
		t.Fatal("expected temperature areas")
	}
}
