// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package mappers

import (
	"sort"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
)

type AreaMapper struct{}

func NewAreaMapper() AreaMapper {
	return AreaMapper{}
}

func (m AreaMapper) ToAreasResponse(document clients.AreaJSONDocument, parent *string) gen.AreasResponse {
	items := make([]gen.Area, 0, len(document.Offices))
	for code, office := range document.Offices {
		if parent != nil && office.Parent != *parent {
			continue
		}

		items = append(items, gen.Area{
			Code:       code,
			Name:       office.Name,
			EnName:     office.EnName,
			OfficeName: office.OfficeName,
			Parent:     office.Parent,
			Children:   append([]string(nil), office.Children...),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})

	return gen.AreasResponse{Items: items}
}

func (m AreaMapper) ToArea(document clients.AreaJSONDocument, areaCode string) (gen.Area, bool) {
	office, ok := document.Offices[areaCode]
	if !ok {
		return gen.Area{}, false
	}

	return gen.Area{
		Code:       areaCode,
		Name:       office.Name,
		EnName:     office.EnName,
		OfficeName: office.OfficeName,
		Parent:     office.Parent,
		Children:   append([]string(nil), office.Children...),
	}, true
}
