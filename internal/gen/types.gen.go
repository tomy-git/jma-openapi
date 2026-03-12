// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

// Code generated for the current OpenAPI contract scaffold. DO NOT EDIT.

package gen

type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type AreasResponse struct {
	Items []Area `json:"items"`
}

type Area struct {
	Code       string   `json:"code"`
	Name       string   `json:"name"`
	EnName     string   `json:"enName"`
	OfficeName string   `json:"officeName"`
	Parent     string   `json:"parent"`
	Children   []string `json:"children"`
}

type ForecastResponse struct {
	Office           OfficeRef               `json:"office"`
	PublishingOffice string                  `json:"publishingOffice"`
	ReportDatetime   string                  `json:"reportDatetime"`
	WeatherAreas     []WeatherArea           `json:"weatherAreas"`
	TemperatureAreas []TemperatureArea       `json:"temperatureAreas"`
}

type OfficeRef struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type WeatherArea struct {
	Code       string                   `json:"code"`
	Name       string                   `json:"name"`
	TimeSeries []WeatherTimeSeriesEntry `json:"timeSeries"`
}

type WeatherTimeSeriesEntry struct {
	Time        string  `json:"time"`
	WeatherCode *string `json:"weatherCode,omitempty"`
	Weather     *string `json:"weather,omitempty"`
	Wind        *string `json:"wind,omitempty"`
	Wave        *string `json:"wave,omitempty"`
	Pop         *string `json:"pop,omitempty"`
}

type TemperatureArea struct {
	Code       string                       `json:"code"`
	Name       string                       `json:"name"`
	TimeSeries []TemperatureTimeSeriesEntry `json:"timeSeries"`
}

type TemperatureTimeSeriesEntry struct {
	Time string  `json:"time"`
	Temp *string `json:"temp,omitempty"`
}

type ErrorEnvelope struct {
	Error ErrorModel `json:"error"`
}

type ErrorModel struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	RequestID string         `json:"requestId"`
	Details   map[string]any `json:"details,omitempty"`
}
