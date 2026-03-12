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

type AreaFilter struct {
	Parent     *string
	Name       *string
	OfficeName *string
	Child      *string
}

func NewAreaMapper() AreaMapper {
	return AreaMapper{}
}

func (m AreaMapper) ToAreasResponse(document clients.AreaJSONDocument, filter AreaFilter) gen.AreasResponse {
	items := make([]gen.Area, 0, len(document.Offices))
	for code, office := range document.Offices {
		if filter.Parent != nil && office.Parent != *filter.Parent {
			continue
		}
		if filter.Name != nil && office.Name != *filter.Name {
			continue
		}
		if filter.OfficeName != nil && office.OfficeName != *filter.OfficeName {
			continue
		}
		if filter.Child != nil && !containsString(office.Children, *filter.Child) {
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

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}

	return false
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
