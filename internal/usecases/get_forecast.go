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
		return gen.ForecastAreaResponse{}, shared.NewAppError(http.StatusNotFound, "FORECAST_AREA_NOT_FOUND", "forecast area code was not found", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, nil)
	}

	return response, nil
}

func (u ForecastUsecase) ListAreas(ctx context.Context, officeCode string) (gen.ForecastAreaListResponse, error) {
	report, office, err := u.loadForecastWithOffice(ctx, officeCode)
	if err != nil {
		return gen.ForecastAreaListResponse{}, err
	}

	_ = office
	return u.forecastMapper.ToForecastAreaListResponse(report)
}

func (u ForecastUsecase) ListWeatherAreas(ctx context.Context, officeCode string) (gen.WeatherAreaListResponse, error) {
	report, _, err := u.loadForecastWithOffice(ctx, officeCode)
	if err != nil {
		return gen.WeatherAreaListResponse{}, err
	}

	return u.forecastMapper.ToWeatherAreaListResponse(report)
}

func (u ForecastUsecase) GetWeatherArea(ctx context.Context, officeCode string, areaCode string) (gen.WeatherAreaResponse, error) {
	report, office, err := u.loadForecastWithOffice(ctx, officeCode)
	if err != nil {
		return gen.WeatherAreaResponse{}, err
	}

	response, ok, err := u.forecastMapper.ToWeatherAreaResponse(report, office, areaCode)
	if err != nil {
		return gen.WeatherAreaResponse{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to normalize weather area payload", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, err)
	}
	if !ok {
		return gen.WeatherAreaResponse{}, shared.NewAppError(http.StatusNotFound, "FORECAST_AREA_NOT_FOUND", "forecast area code was not found", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, nil)
	}

	return response, nil
}

func (u ForecastUsecase) ListTemperatureAreas(ctx context.Context, officeCode string) (gen.TemperatureAreaListResponse, error) {
	report, _, err := u.loadForecastWithOffice(ctx, officeCode)
	if err != nil {
		return gen.TemperatureAreaListResponse{}, err
	}

	return u.forecastMapper.ToTemperatureAreaListResponse(report)
}

func (u ForecastUsecase) GetTemperatureArea(ctx context.Context, officeCode string, areaCode string) (gen.TemperatureAreaResponse, error) {
	report, office, err := u.loadForecastWithOffice(ctx, officeCode)
	if err != nil {
		return gen.TemperatureAreaResponse{}, err
	}

	response, ok, err := u.forecastMapper.ToTemperatureAreaResponse(report, office, areaCode)
	if err != nil {
		return gen.TemperatureAreaResponse{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to normalize temperature area payload", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, err)
	}
	if !ok {
		return gen.TemperatureAreaResponse{}, shared.NewAppError(http.StatusNotFound, "FORECAST_AREA_NOT_FOUND", "forecast area code was not found", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, nil)
	}

	return response, nil
}

func (u ForecastUsecase) ResolveAreas(ctx context.Context, officeCode string, params gen.GetV1ForecastsOfficeCodeAreasResolveParams) (gen.ForecastAreaListResponse, error) {
	report, _, err := u.loadForecastWithOffice(ctx, officeCode)
	if err != nil {
		return gen.ForecastAreaListResponse{}, err
	}

	matchMode := gen.AreaMatchMode("exact")
	if params.MatchMode != nil {
		matchMode = *params.MatchMode
	}
	switch matchMode {
	case gen.AreaMatchMode("exact"), gen.AreaMatchMode("prefix"), gen.AreaMatchMode("partial"), gen.AreaMatchMode("suggested"):
	default:
		return gen.ForecastAreaListResponse{}, shared.NewAppError(http.StatusBadRequest, "INVALID_MATCH_MODE", "match mode was invalid", map[string]any{"matchMode": matchMode}, nil)
	}

	kind := gen.ForecastAreaKind("any")
	if params.Kind != nil {
		kind = *params.Kind
	}

	response, err := u.forecastMapper.ResolveForecastAreas(report, params.Q, kind, matchMode)
	if err != nil {
		return gen.ForecastAreaListResponse{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to resolve forecast areas", map[string]any{"officeCode": officeCode}, err)
	}

	return response, nil
}

func (u ForecastUsecase) GetAreaLatest(ctx context.Context, officeCode string, areaCode string) (gen.ForecastAreaResponse, error) {
	report, office, err := u.loadForecastWithOffice(ctx, officeCode)
	if err != nil {
		return gen.ForecastAreaResponse{}, err
	}

	response, ok, err := u.forecastMapper.ToForecastAreaLatestResponse(report, office, areaCode)
	if err != nil {
		return gen.ForecastAreaResponse{}, shared.NewAppError(http.StatusBadGateway, "UPSTREAM_SCHEMA_MISMATCH", "failed to normalize latest forecast area payload", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, err)
	}
	if !ok {
		return gen.ForecastAreaResponse{}, shared.NewAppError(http.StatusNotFound, "FORECAST_AREA_NOT_FOUND", "forecast area code was not found", map[string]any{"officeCode": officeCode, "areaCode": areaCode}, nil)
	}

	return response, nil
}

func (u ForecastUsecase) GetAreaTimeseries(ctx context.Context, officeCode string, areaCode string) (gen.ForecastAreaResponse, error) {
	return u.GetArea(ctx, officeCode, areaCode)
}

func (u ForecastUsecase) loadForecastWithOffice(ctx context.Context, officeCode string) (clients.ForecastReportJSON, gen.Area, error) {
	areas, err := u.client.FetchAreaDocument(ctx)
	if err != nil {
		return clients.ForecastReportJSON{}, gen.Area{}, err
	}

	office, ok := u.areaMapper.ToArea(areas, officeCode)
	if !ok {
		return clients.ForecastReportJSON{}, gen.Area{}, shared.NewAppError(http.StatusNotFound, "OFFICE_NOT_FOUND", "office code was not found", map[string]any{"officeCode": officeCode}, nil)
	}

	report, err := u.client.FetchForecastDocument(ctx, officeCode)
	if err != nil {
		return clients.ForecastReportJSON{}, gen.Area{}, err
	}

	return report, office, nil
}
