// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package usecases

import (
	"context"
	"net/http"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
	"github.com/tomy-git/jma-openapi/internal/mappers"
	"github.com/tomy-git/jma-openapi/internal/shared"
)

type ForecastUsecase struct {
	client         clients.JMAClient
	areaMapper     mappers.AreaMapper
	forecastMapper mappers.ForecastMapper
}

func NewForecastUsecase(client clients.JMAClient, areaMapper mappers.AreaMapper, forecastMapper mappers.ForecastMapper) ForecastUsecase {
	return ForecastUsecase{
		client:         client,
		areaMapper:     areaMapper,
		forecastMapper: forecastMapper,
	}
}

func (u ForecastUsecase) Get(ctx context.Context, officeCode string) (gen.ForecastResponse, error) {
	areas, err := u.client.FetchAreaDocument(ctx)
	if err != nil {
		return gen.ForecastResponse{}, err
	}

	office, ok := u.areaMapper.ToArea(areas, officeCode)
	if !ok {
		return gen.ForecastResponse{}, shared.NewAppError(http.StatusNotFound, "OFFICE_NOT_FOUND", "office code was not found", map[string]any{"officeCode": officeCode}, nil)
	}

	report, err := u.client.FetchForecastDocument(ctx, officeCode)
	if err != nil {
		return gen.ForecastResponse{}, err
	}

	response, err := u.forecastMapper.ToForecastResponse(report, office)
	if err != nil {
		return gen.ForecastResponse{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to map forecast response", map[string]any{"officeCode": officeCode}, err)
	}

	return response, nil
}

func (u ForecastUsecase) GetArea(ctx context.Context, officeCode string, areaCode string) (gen.ForecastAreaResponse, error) {
	areas, err := u.client.FetchAreaDocument(ctx)
	if err != nil {
		return gen.ForecastAreaResponse{}, err
	}

	office, ok := u.areaMapper.ToArea(areas, officeCode)
	if !ok {
		return gen.ForecastAreaResponse{}, shared.NewAppError(http.StatusNotFound, "OFFICE_NOT_FOUND", "office code was not found", map[string]any{"officeCode": officeCode}, nil)
	}

	report, err := u.client.FetchWeatherAreaForecastDocument(ctx, officeCode)
	if err != nil {
		return gen.ForecastAreaResponse{}, err
	}

	response, ok, err := u.forecastMapper.ToForecastAreaResponse(report, office, areaCode)
	if err != nil {
		return gen.ForecastAreaResponse{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to normalize forecast payload", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, err)
	}
	if !ok {
		return gen.ForecastAreaResponse{}, shared.NewAppError(http.StatusNotFound, "WEATHER_AREA_NOT_FOUND", "weather area code was not found", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, nil)
	}

	return response, nil
}
