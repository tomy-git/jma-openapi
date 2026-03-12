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

type AreaUsecase struct {
	client clients.JMAClient
	mapper mappers.AreaMapper
}

func NewAreaUsecase(client clients.JMAClient, mapper mappers.AreaMapper) AreaUsecase {
	return AreaUsecase{
		client: client,
		mapper: mapper,
	}
}

func (u AreaUsecase) List(ctx context.Context, filter mappers.AreaFilter) (gen.AreasResponse, error) {
	switch filter.NameMode {
	case "", gen.AreaMatchMode("exact"), gen.AreaMatchMode("prefix"), gen.AreaMatchMode("partial"), gen.AreaMatchMode("suggested"):
	default:
		return gen.AreasResponse{}, shared.NewAppError(http.StatusBadRequest, "INVALID_MATCH_MODE", "name match mode was invalid", map[string]any{"nameMatchMode": filter.NameMode}, nil)
	}

	document, err := u.client.FetchAreaDocument(ctx)
	if err != nil {
		return gen.AreasResponse{}, err
	}

	return u.mapper.ToAreasResponse(document, filter), nil
}

func (u AreaUsecase) Get(ctx context.Context, areaCode string) (gen.Area, error) {
	document, err := u.client.FetchAreaDocument(ctx)
	if err != nil {
		return gen.Area{}, err
	}

	area, ok := u.mapper.ToArea(document, areaCode)
	if !ok {
		return gen.Area{}, shared.NewAppError(http.StatusNotFound, "AREA_NOT_FOUND", "area code was not found", map[string]any{"areaCode": areaCode}, nil)
	}

	return area, nil
}
