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

	return u.forecastMapper.ToForecastResponse(report, office), nil
}
