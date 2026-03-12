// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package mappers

import (
	"sort"
	"strings"
	"time"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
)

type ForecastMapper struct{}

func NewForecastMapper() ForecastMapper {
	return ForecastMapper{}
}

func (m ForecastMapper) ToForecastResponse(report clients.ForecastReportJSON, office gen.Area) (gen.ForecastResponse, error) {
	reportDatetime, err := parseForecastTime(report.ReportDatetime)
	if err != nil {
		return gen.ForecastResponse{}, err
	}

	weatherAreas, err := m.mergeWeatherAreas(report)
	if err != nil {
		return gen.ForecastResponse{}, err
	}

	temperatureAreas, err := m.temperatureAreas(report)
	if err != nil {
		return gen.ForecastResponse{}, err
	}

	response := gen.ForecastResponse{
		Office: gen.OfficeRef{
			Code: office.Code,
			Name: office.Name,
		},
		PublishingOffice: report.PublishingOffice,
		ReportDatetime:   reportDatetime,
		WeatherAreas:     weatherAreas,
		TemperatureAreas: temperatureAreas,
	}

	return response, nil
}

func (m ForecastMapper) ToForecastAreaListResponse(report clients.ForecastReportJSON) (gen.ForecastAreaListResponse, error) {
	weatherItems, err := m.ToWeatherAreaListResponse(report)
	if err != nil {
		return gen.ForecastAreaListResponse{}, err
	}
	temperatureItems, err := m.ToTemperatureAreaListResponse(report)
	if err != nil {
		return gen.ForecastAreaListResponse{}, err
	}

	items := append(append([]gen.ForecastAreaRef(nil), weatherItems.Items...), temperatureItems.Items...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Kind == items[j].Kind {
			return items[i].Code < items[j].Code
		}

		return items[i].Kind < items[j].Kind
	})

	return gen.ForecastAreaListResponse{Items: items}, nil
}

func (m ForecastMapper) ToWeatherAreaListResponse(report clients.ForecastReportJSON) (gen.WeatherAreaListResponse, error) {
	weatherAreas, err := m.mergeWeatherAreas(report)
	if err != nil {
		return gen.WeatherAreaListResponse{}, err
	}

	items := make([]gen.ForecastAreaRef, 0, len(weatherAreas))
	for _, area := range weatherAreas {
		items = append(items, gen.ForecastAreaRef{Kind: gen.ForecastAreaKind("weather"), Code: area.Code, Name: area.Name})
	}

	return gen.WeatherAreaListResponse{Items: items}, nil
}

func (m ForecastMapper) ToTemperatureAreaListResponse(report clients.ForecastReportJSON) (gen.TemperatureAreaListResponse, error) {
	temperatureAreas, err := m.temperatureAreas(report)
	if err != nil {
		return gen.TemperatureAreaListResponse{}, err
	}

	items := make([]gen.ForecastAreaRef, 0, len(temperatureAreas))
	for _, area := range temperatureAreas {
		items = append(items, gen.ForecastAreaRef{Kind: gen.ForecastAreaKind("temperature"), Code: area.Code, Name: area.Name})
	}

	return gen.TemperatureAreaListResponse{Items: items}, nil
}

func (m ForecastMapper) ToForecastAreaResponse(report clients.ForecastReportJSON, office gen.Area, areaCode string) (gen.ForecastAreaResponse, bool, error) {
	reportDatetime, err := parseForecastTime(report.ReportDatetime)
	if err != nil {
		return gen.ForecastAreaResponse{}, false, err
	}

	weatherAreas, err := m.mergeWeatherAreas(report)
	if err != nil {
		return gen.ForecastAreaResponse{}, false, err
	}
	temperatureAreas, err := m.temperatureAreas(report)
	if err != nil {
		return gen.ForecastAreaResponse{}, false, err
	}

	for _, weatherArea := range weatherAreas {
		if weatherArea.Code != areaCode {
			continue
		}

		return gen.ForecastAreaResponse{
			Office: gen.OfficeRef{
				Code: office.Code,
				Name: office.Name,
			},
			PublishingOffice: report.PublishingOffice,
			ReportDatetime:   reportDatetime,
			WeatherArea:      &weatherArea,
		}, true, nil
	}

	for _, temperatureArea := range temperatureAreas {
		if temperatureArea.Code != areaCode {
			continue
		}

		return gen.ForecastAreaResponse{
			Office: gen.OfficeRef{
				Code: office.Code,
				Name: office.Name,
			},
			PublishingOffice: report.PublishingOffice,
			ReportDatetime:   reportDatetime,
			TemperatureArea:  &temperatureArea,
		}, true, nil
	}

	return gen.ForecastAreaResponse{}, false, nil
}

func (m ForecastMapper) ToWeatherAreaResponse(report clients.ForecastReportJSON, office gen.Area, areaCode string) (gen.WeatherAreaResponse, bool, error) {
	reportDatetime, err := parseForecastTime(report.ReportDatetime)
	if err != nil {
		return gen.WeatherAreaResponse{}, false, err
	}

	weatherAreas, err := m.mergeWeatherAreas(report)
	if err != nil {
		return gen.WeatherAreaResponse{}, false, err
	}

	for _, area := range weatherAreas {
		if area.Code == areaCode {
			return gen.WeatherAreaResponse{
				Office:           gen.OfficeRef{Code: office.Code, Name: office.Name},
				PublishingOffice: report.PublishingOffice,
				ReportDatetime:   reportDatetime,
				WeatherArea:      area,
			}, true, nil
		}
	}

	return gen.WeatherAreaResponse{}, false, nil
}

func (m ForecastMapper) ToTemperatureAreaResponse(report clients.ForecastReportJSON, office gen.Area, areaCode string) (gen.TemperatureAreaResponse, bool, error) {
	reportDatetime, err := parseForecastTime(report.ReportDatetime)
	if err != nil {
		return gen.TemperatureAreaResponse{}, false, err
	}

	temperatureAreas, err := m.temperatureAreas(report)
	if err != nil {
		return gen.TemperatureAreaResponse{}, false, err
	}

	for _, area := range temperatureAreas {
		if area.Code == areaCode {
			return gen.TemperatureAreaResponse{
				Office:           gen.OfficeRef{Code: office.Code, Name: office.Name},
				PublishingOffice: report.PublishingOffice,
				ReportDatetime:   reportDatetime,
				TemperatureArea:  area,
			}, true, nil
		}
	}

	return gen.TemperatureAreaResponse{}, false, nil
}

func (m ForecastMapper) ResolveForecastAreas(report clients.ForecastReportJSON, query string, kind gen.ForecastAreaKind, mode gen.AreaMatchMode) (gen.ForecastAreaListResponse, error) {
	query = normalizeSearchText(query)
	if query == "" {
		return gen.ForecastAreaListResponse{Items: []gen.ForecastAreaRef{}}, nil
	}

	items := make([]scoredAreaRef, 0)
	if kind == "" || kind == gen.ForecastAreaKind("any") || kind == gen.ForecastAreaKind("weather") {
		weatherResponse, err := m.ToWeatherAreaListResponse(report)
		if err != nil {
			return gen.ForecastAreaListResponse{}, err
		}
		items = append(items, matchAreaRefs(weatherResponse.Items, query, mode)...)
	}
	if kind == "" || kind == gen.ForecastAreaKind("any") || kind == gen.ForecastAreaKind("temperature") {
		temperatureResponse, err := m.ToTemperatureAreaListResponse(report)
		if err != nil {
			return gen.ForecastAreaListResponse{}, err
		}
		items = append(items, matchAreaRefs(temperatureResponse.Items, query, mode)...)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].score == items[j].score {
			if items[i].ref.Kind == items[j].ref.Kind {
				return items[i].ref.Code < items[j].ref.Code
			}
			return items[i].ref.Kind < items[j].ref.Kind
		}
		return items[i].score > items[j].score
	})

	refs := make([]gen.ForecastAreaRef, 0, len(items))
	for _, item := range items {
		refs = append(refs, item.ref)
	}

	return gen.ForecastAreaListResponse{Items: refs}, nil
}

func (m ForecastMapper) ToForecastAreaLatestResponse(report clients.ForecastReportJSON, office gen.Area, areaCode string) (gen.ForecastAreaResponse, bool, error) {
	response, ok, err := m.ToForecastAreaResponse(report, office, areaCode)
	if !ok || err != nil {
		return response, ok, err
	}

	if response.WeatherArea != nil && len(response.WeatherArea.TimeSeries) > 0 {
		response.WeatherArea.TimeSeries = response.WeatherArea.TimeSeries[:1]
	}
	if response.TemperatureArea != nil && len(response.TemperatureArea.TimeSeries) > 0 {
		response.TemperatureArea.TimeSeries = response.TemperatureArea.TimeSeries[:1]
	}

	return response, true, nil
}

type scoredAreaRef struct {
	ref   gen.ForecastAreaRef
	score int
}

func matchAreaRefs(items []gen.ForecastAreaRef, query string, mode gen.AreaMatchMode) []scoredAreaRef {
	matched := make([]scoredAreaRef, 0, len(items))
	for _, item := range items {
		score, ok := matchSearchMode(normalizeSearchText(item.Name), query, mode)
		if !ok {
			continue
		}
		matched = append(matched, scoredAreaRef{ref: item, score: score})
	}

	return matched
}

func matchSearchMode(name string, query string, mode gen.AreaMatchMode) (int, bool) {
	switch mode {
	case "", gen.AreaMatchMode("exact"):
		return 100, name == query
	case gen.AreaMatchMode("prefix"):
		return 80, strings.HasPrefix(name, query)
	case gen.AreaMatchMode("partial"):
		return 60, strings.Contains(name, query)
	case gen.AreaMatchMode("suggested"):
		switch {
		case name == query:
			return 100, true
		case strings.HasPrefix(name, query):
			return 80, true
		case strings.Contains(name, query):
			return 60, true
		default:
			return 0, false
		}
	default:
		return 0, false
	}
}

func normalizeSearchText(value string) string {
	return strings.TrimSpace(value)
}

func (m ForecastMapper) mergeWeatherAreas(report clients.ForecastReportJSON) ([]gen.WeatherArea, error) {
	if len(report.WeatherSeries) == 0 {
		return nil, nil
	}

	popLookup := map[string]clients.ForecastAreaPopJSON{}
	popTimes := []string(nil)
	if len(report.PopSeries) > 0 {
		popTimes = report.PopSeries[0].TimeDefines
		for _, area := range report.PopSeries[0].Areas {
			popLookup[area.Area.Code] = area
		}
	}

	items := make([]gen.WeatherArea, 0, len(report.WeatherSeries[0].Areas))
	for _, area := range report.WeatherSeries[0].Areas {
		entries := make([]gen.WeatherTimeSeriesEntry, 0, len(report.WeatherSeries[0].TimeDefines))
		for idx, timeValue := range report.WeatherSeries[0].TimeDefines {
			parsedTime, err := parseForecastTime(timeValue)
			if err != nil {
				return nil, err
			}

			entry := gen.WeatherTimeSeriesEntry{
				Time:        parsedTime,
				WeatherCode: stringAt(area.WeatherCodes, idx),
				Weather:     stringAt(area.Weathers, idx),
				Wind:        stringAt(area.Winds, idx),
				Wave:        stringAt(area.Waves, idx),
			}

			if popArea, ok := popLookup[area.Area.Code]; ok {
				if popIndex := indexOf(popTimes, timeValue); popIndex >= 0 {
					entry.Pop = stringAt(popArea.Pops, popIndex)
				}
			}

			entries = append(entries, entry)
		}

		items = append(items, gen.WeatherArea{
			Code:       area.Area.Code,
			Name:       area.Area.Name,
			TimeSeries: entries,
		})
	}

	return items, nil
}

func (m ForecastMapper) temperatureAreas(report clients.ForecastReportJSON) ([]gen.TemperatureArea, error) {
	if len(report.TempSeries) == 0 {
		return nil, nil
	}

	items := make([]gen.TemperatureArea, 0, len(report.TempSeries[0].Areas))
	for _, area := range report.TempSeries[0].Areas {
		entries := make([]gen.TemperatureTimeSeriesEntry, 0, len(report.TempSeries[0].TimeDefines))
		for idx, timeValue := range report.TempSeries[0].TimeDefines {
			parsedTime, err := parseForecastTime(timeValue)
			if err != nil {
				return nil, err
			}

			entries = append(entries, gen.TemperatureTimeSeriesEntry{
				Time: parsedTime,
				Temp: stringAt(area.Temps, idx),
			})
		}

		items = append(items, gen.TemperatureArea{
			Code:       area.Area.Code,
			Name:       area.Area.Name,
			TimeSeries: entries,
		})
	}

	return items, nil
}

func parseForecastTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func stringAt(items []string, index int) *string {
	if index < 0 || index >= len(items) {
		return nil
	}

	value := items[index]
	return &value
}

func indexOf(items []string, target string) int {
	for idx, item := range items {
		if item == target {
			return idx
		}
	}

	return -1
}
