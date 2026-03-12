// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package mappers

import (
	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
)

type ForecastMapper struct{}

func NewForecastMapper() ForecastMapper {
	return ForecastMapper{}
}

func (m ForecastMapper) ToForecastResponse(report clients.ForecastReportJSON, office gen.Area) gen.ForecastResponse {
	response := gen.ForecastResponse{
		Office: gen.OfficeRef{
			Code: office.Code,
			Name: office.Name,
		},
		PublishingOffice: report.PublishingOffice,
		ReportDatetime:   report.ReportDatetime,
		WeatherAreas:     m.mergeWeatherAreas(report),
		TemperatureAreas: m.temperatureAreas(report),
	}

	return response
}

func (m ForecastMapper) mergeWeatherAreas(report clients.ForecastReportJSON) []gen.WeatherArea {
	if len(report.WeatherSeries) == 0 {
		return nil
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
			entry := gen.WeatherTimeSeriesEntry{
				Time:        timeValue,
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

	return items
}

func (m ForecastMapper) temperatureAreas(report clients.ForecastReportJSON) []gen.TemperatureArea {
	if len(report.TempSeries) == 0 {
		return nil
	}

	items := make([]gen.TemperatureArea, 0, len(report.TempSeries[0].Areas))
	for _, area := range report.TempSeries[0].Areas {
		entries := make([]gen.TemperatureTimeSeriesEntry, 0, len(report.TempSeries[0].TimeDefines))
		for idx, timeValue := range report.TempSeries[0].TimeDefines {
			entries = append(entries, gen.TemperatureTimeSeriesEntry{
				Time: timeValue,
				Temp: stringAt(area.Temps, idx),
			})
		}

		items = append(items, gen.TemperatureArea{
			Code:       area.Area.Code,
			Name:       area.Area.Name,
			TimeSeries: entries,
		})
	}

	return items
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
