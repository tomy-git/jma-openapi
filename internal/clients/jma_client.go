// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tomy-git/jma-openapi/internal/shared"
)

type AreaOfficeJSON struct {
	Name       string   `json:"name"`
	EnName     string   `json:"enName"`
	OfficeName string   `json:"officeName"`
	Parent     string   `json:"parent"`
	Children   []string `json:"children"`
}

type AreaJSONDocument struct {
	Offices map[string]AreaOfficeJSON `json:"offices"`
}

type ForecastAreaReference struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ForecastAreaWeatherJSON struct {
	Area         ForecastAreaReference `json:"area"`
	WeatherCodes []string              `json:"weatherCodes"`
	Weathers     []string              `json:"weathers"`
	Winds        []string              `json:"winds"`
	Waves        []string              `json:"waves"`
}

type ForecastAreaPopJSON struct {
	Area ForecastAreaReference `json:"area"`
	Pops []string              `json:"pops"`
}

type ForecastAreaTempJSON struct {
	Area  ForecastAreaReference `json:"area"`
	Temps []string              `json:"temps"`
}

type ForecastTimeSeriesWeatherJSON struct {
	TimeDefines []string                  `json:"timeDefines"`
	Areas       []ForecastAreaWeatherJSON `json:"areas"`
}

type ForecastTimeSeriesPopJSON struct {
	TimeDefines []string              `json:"timeDefines"`
	Areas       []ForecastAreaPopJSON `json:"areas"`
}

type ForecastTimeSeriesTempJSON struct {
	TimeDefines []string               `json:"timeDefines"`
	Areas       []ForecastAreaTempJSON `json:"areas"`
}

type ForecastReportJSON struct {
	PublishingOffice string `json:"publishingOffice"`
	ReportDatetime   string `json:"reportDatetime"`
	WeatherSeries    []ForecastTimeSeriesWeatherJSON
	PopSeries        []ForecastTimeSeriesPopJSON
	TempSeries       []ForecastTimeSeriesTempJSON
}

type forecastReportWire struct {
	PublishingOffice string            `json:"publishingOffice"`
	ReportDatetime   string            `json:"reportDatetime"`
	TimeSeries       []json.RawMessage `json:"timeSeries"`
}

type JMAClient interface {
	FetchAreaDocument(ctx context.Context) (AreaJSONDocument, error)
	FetchForecastDocument(ctx context.Context, officeCode string) (ForecastReportJSON, error)
	FetchWeatherAreaForecastDocument(ctx context.Context, officeCode string) (ForecastReportJSON, error)
}

type HTTPJMAClient struct {
	baseURL    string
	httpClient *http.Client

	mu         sync.RWMutex
	cachedArea cachedAreaDocument
}

type cachedAreaDocument struct {
	expiresAt time.Time
	document  AreaJSONDocument
	ok        bool
}

func NewHTTPJMAClient(baseURL string) *HTTPJMAClient {
	return &HTTPJMAClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPJMAClient) FetchAreaDocument(ctx context.Context) (AreaJSONDocument, error) {
	c.mu.RLock()
	if c.cachedArea.ok && time.Now().Before(c.cachedArea.expiresAt) {
		document := c.cachedArea.document
		c.mu.RUnlock()
		return document, nil
	}
	c.mu.RUnlock()

	url := fmt.Sprintf("%s/common/const/area.json", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return AreaJSONDocument{}, shared.NewAppError(http.StatusInternalServerError, "REQUEST_BUILD_FAILED", "failed to build upstream request", nil, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AreaJSONDocument{}, shared.NewAppError(http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", "failed to fetch area metadata from upstream", nil, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return AreaJSONDocument{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_BAD_RESPONSE", "upstream returned non-200 for area metadata", map[string]any{"status": resp.StatusCode}, nil)
	}

	var document AreaJSONDocument
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return AreaJSONDocument{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_DECODE_FAILED", "failed to decode area metadata", nil, err)
	}

	ttl := cacheTTL(resp.Header.Get("Cache-Control"))
	c.mu.Lock()
	c.cachedArea = cachedAreaDocument{
		expiresAt: time.Now().Add(ttl),
		document:  document,
		ok:        true,
	}
	c.mu.Unlock()

	return document, nil
}

func (c *HTTPJMAClient) FetchForecastDocument(ctx context.Context, officeCode string) (ForecastReportJSON, error) {
	reportWire, err := c.fetchForecastReportWire(ctx, officeCode)
	if err != nil {
		return ForecastReportJSON{}, err
	}

	if len(reportWire.TimeSeries) < 3 {
		return ForecastReportJSON{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "forecast payload did not include required timeSeries entries", map[string]any{"officeCode": officeCode}, nil)
	}

	return decodeForecastReport(reportWire, officeCode, true)
}

func (c *HTTPJMAClient) FetchWeatherAreaForecastDocument(ctx context.Context, officeCode string) (ForecastReportJSON, error) {
	reportWire, err := c.fetchForecastReportWire(ctx, officeCode)
	if err != nil {
		return ForecastReportJSON{}, err
	}

	if len(reportWire.TimeSeries) == 0 {
		return ForecastReportJSON{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "forecast payload did not include required weather timeSeries", map[string]any{"officeCode": officeCode}, nil)
	}

	return decodeForecastReport(reportWire, officeCode, false)
}

func (c *HTTPJMAClient) fetchForecastReportWire(ctx context.Context, officeCode string) (forecastReportWire, error) {
	url := fmt.Sprintf("%s/forecast/data/forecast/%s.json", c.baseURL, officeCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return forecastReportWire{}, shared.NewAppError(http.StatusInternalServerError, "REQUEST_BUILD_FAILED", "failed to build upstream request", nil, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return forecastReportWire{}, shared.NewAppError(http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", "failed to fetch forecast from upstream", map[string]any{"officeCode": officeCode}, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return forecastReportWire{}, shared.NewAppError(http.StatusNotFound, "OFFICE_NOT_FOUND", "office code was not found", map[string]any{"officeCode": officeCode}, nil)
	}

	if resp.StatusCode != http.StatusOK {
		return forecastReportWire{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_BAD_RESPONSE", "upstream returned non-200 for forecast", map[string]any{"status": resp.StatusCode, "officeCode": officeCode}, nil)
	}

	var reports []forecastReportWire
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return forecastReportWire{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_DECODE_FAILED", "failed to decode forecast", map[string]any{"officeCode": officeCode}, err)
	}

	if len(reports) == 0 {
		return forecastReportWire{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_EMPTY_RESPONSE", "forecast payload was empty", map[string]any{"officeCode": officeCode}, nil)
	}

	return reports[0], nil
}

func decodeForecastReport(reportWire forecastReportWire, officeCode string, requireTempSeries bool) (ForecastReportJSON, error) {
	report := ForecastReportJSON{
		PublishingOffice: reportWire.PublishingOffice,
		ReportDatetime:   reportWire.ReportDatetime,
	}

	if len(reportWire.TimeSeries) > 0 {
		var weather ForecastTimeSeriesWeatherJSON
		if err := json.Unmarshal(reportWire.TimeSeries[0], &weather); err != nil {
			return ForecastReportJSON{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to decode weather forecast timeSeries", map[string]any{"officeCode": officeCode, "index": 0}, err)
		}
		report.WeatherSeries = []ForecastTimeSeriesWeatherJSON{weather}
	}
	if len(reportWire.TimeSeries) > 1 {
		var pops ForecastTimeSeriesPopJSON
		if err := json.Unmarshal(reportWire.TimeSeries[1], &pops); err != nil {
			return ForecastReportJSON{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to decode precipitation probability timeSeries", map[string]any{"officeCode": officeCode, "index": 1}, err)
		}
		report.PopSeries = []ForecastTimeSeriesPopJSON{pops}
	}
	if len(reportWire.TimeSeries) > 2 {
		var temps ForecastTimeSeriesTempJSON
		if err := json.Unmarshal(reportWire.TimeSeries[2], &temps); err != nil {
			return ForecastReportJSON{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to decode temperature forecast timeSeries", map[string]any{"officeCode": officeCode, "index": 2}, err)
		}
		report.TempSeries = []ForecastTimeSeriesTempJSON{temps}
	} else if requireTempSeries {
		return ForecastReportJSON{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "forecast payload did not include required timeSeries entries", map[string]any{"officeCode": officeCode}, nil)
	}

	return report, nil
}

func cacheTTL(cacheControl string) time.Duration {
	for _, part := range strings.Split(cacheControl, ",") {
		part = strings.TrimSpace(part)
		if !strings.HasPrefix(part, "max-age=") {
			continue
		}

		seconds, err := strconv.Atoi(strings.TrimPrefix(part, "max-age="))
		if err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}

	return time.Minute
}
