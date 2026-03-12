// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package mappers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	response, err := mapper.ToForecastResponse(report, gen.Area{Code: "130000", Name: "東京都"})
	if err != nil {
		t.Fatal(err)
	}

	if response.Office.Code != "130000" {
		t.Fatalf("expected office code 130000, got %s", response.Office.Code)
	}
	if !response.ReportDatetime.Equal(mustParseTime(t, "2026-03-12T05:00:00+09:00")) {
		t.Fatalf("expected report datetime to be parsed, got %s", response.ReportDatetime)
	}
	if len(response.WeatherAreas) == 0 {
		t.Fatal("expected weather areas")
	}
	if len(response.TemperatureAreas) == 0 {
		t.Fatal("expected temperature areas")
	}
	if !response.WeatherAreas[0].TimeSeries[0].Time.Equal(mustParseTime(t, "2026-03-12T05:00:00+09:00")) {
		t.Fatalf("expected weather time to be parsed, got %s", response.WeatherAreas[0].TimeSeries[0].Time)
	}
}

func TestForecastMapper_ToForecastAreaResponse(t *testing.T) {
	t.Parallel()

	report := loadForecastReport(t)

	mapper := NewForecastMapper()
	response, ok, err := mapper.ToForecastAreaResponse(report, gen.Area{Code: "130000", Name: "東京都"}, "130010")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected target weather area")
	}
	if response.Office.Code != "130000" {
		t.Fatalf("expected office code 130000, got %s", response.Office.Code)
	}
	if response.WeatherArea == nil {
		t.Fatal("expected weather area")
	}
	if response.WeatherArea.Code != "130010" {
		t.Fatalf("expected weather area code 130010, got %s", response.WeatherArea.Code)
	}
	if len(response.WeatherArea.TimeSeries) == 0 {
		t.Fatal("expected weather area time series")
	}
	if response.TemperatureArea != nil {
		t.Fatal("did not expect temperature area for weather area lookup")
	}

	response, ok, err = mapper.ToForecastAreaResponse(report, gen.Area{Code: "130000", Name: "東京都"}, "44132")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected target temperature area")
	}
	if response.TemperatureArea == nil {
		t.Fatal("expected temperature area")
	}
	if response.TemperatureArea.Code != "44132" {
		t.Fatalf("expected temperature area code 44132, got %s", response.TemperatureArea.Code)
	}
	if response.WeatherArea != nil {
		t.Fatal("did not expect weather area for temperature area lookup")
	}
	if _, ok, err := mapper.ToForecastAreaResponse(report, gen.Area{Code: "130000", Name: "東京都"}, "999999"); err != nil || ok {
		t.Fatalf("expected area not found without error, got ok=%v err=%v", ok, err)
	}
}

func TestForecastMapper_AreaListsAndResolve(t *testing.T) {
	t.Parallel()

	report := loadForecastReport(t)
	mapper := NewForecastMapper()

	allAreas, err := mapper.ToForecastAreaListResponse(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(allAreas.Items) != 8 {
		t.Fatalf("expected 8 area refs, got %d", len(allAreas.Items))
	}

	weatherAreas, err := mapper.ToWeatherAreaListResponse(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(weatherAreas.Items) != 4 || weatherAreas.Items[0].Kind != gen.ForecastAreaKind("weather") {
		t.Fatalf("unexpected weather areas: %+v", weatherAreas.Items)
	}

	resolved, err := mapper.ResolveForecastAreas(report, "東京", gen.ForecastAreaKind("temperature"), gen.AreaMatchMode("exact"))
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved.Items) != 1 || resolved.Items[0].Code != "44132" {
		t.Fatalf("expected Tokyo temperature area, got %+v", resolved.Items)
	}
}

func TestForecastMapper_ToForecastAreaLatestResponse(t *testing.T) {
	t.Parallel()

	report := loadForecastReport(t)
	mapper := NewForecastMapper()

	response, ok, err := mapper.ToForecastAreaLatestResponse(report, gen.Area{Code: "130000", Name: "東京都"}, "130010")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected target weather area")
	}
	if response.WeatherArea == nil || len(response.WeatherArea.TimeSeries) != 1 {
		t.Fatalf("expected a single weather timeseries entry, got %+v", response.WeatherArea)
	}
}

func TestForecastMapper_AreaListsAndResolveAndLatest(t *testing.T) {
	t.Parallel()

	report := loadForecastReport(t)
	mapper := NewForecastMapper()

	allAreas, err := mapper.ToForecastAreaListResponse(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(allAreas.Items) != 8 {
		t.Fatalf("expected 8 area refs, got %d", len(allAreas.Items))
	}

	weatherAreas, err := mapper.ToWeatherAreaListResponse(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(weatherAreas.Items) != 4 || weatherAreas.Items[0].Kind != gen.Weather {
		t.Fatalf("expected weather area refs, got %+v", weatherAreas.Items)
	}

	temperatureAreas, err := mapper.ToTemperatureAreaListResponse(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(temperatureAreas.Items) != 4 || temperatureAreas.Items[0].Kind != gen.Temperature {
		t.Fatalf("expected temperature area refs, got %+v", temperatureAreas.Items)
	}

	resolved, err := mapper.ResolveForecastAreas(report, "東京", gen.Any, gen.Suggested)
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved.Items) == 0 || resolved.Items[0].Code != "44132" {
		t.Fatalf("expected Tokyo weather area first, got %+v", resolved.Items)
	}

	latest, ok, err := mapper.ToForecastAreaLatestResponse(report, gen.Area{Code: "130000", Name: "東京都"}, "44132")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || latest.TemperatureArea == nil {
		t.Fatal("expected latest temperature area")
	}
	if len(latest.TemperatureArea.TimeSeries) != 1 {
		t.Fatalf("expected single latest entry, got %d", len(latest.TemperatureArea.TimeSeries))
	}
}

func loadForecastReport(t *testing.T) clients.ForecastReportJSON {
	t.Helper()

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

	return report
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}

	return parsed
}
