// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package mappers

import (
	"sort"
	"strings"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
)

type AreaMapper struct{}

type AreaFilter struct {
	Parent     *string
	Name       *string
	NameMode   gen.AreaMatchMode
	OfficeName *string
	Child      *string
}

func NewAreaMapper() AreaMapper {
	return AreaMapper{}
}

func (m AreaMapper) ToAreasResponse(document clients.AreaJSONDocument, filter AreaFilter) gen.AreasResponse {
	type scoredArea struct {
		area  gen.Area
		score int
	}

	items := make([]scoredArea, 0, len(document.Offices))
	for code, office := range document.Offices {
		if filter.Parent != nil && office.Parent != *filter.Parent {
			continue
		}
		nameScore, ok := matchAreaName(filter.NameMode, office.Name, filter.Name)
		if !ok {
			continue
		}
		if filter.OfficeName != nil && office.OfficeName != *filter.OfficeName {
			continue
		}
		if filter.Child != nil && !containsString(office.Children, *filter.Child) {
			continue
		}

		items = append(items, scoredArea{
			score: nameScore,
			area: gen.Area{
				Code:       code,
				Name:       office.Name,
				EnName:     office.EnName,
				OfficeName: office.OfficeName,
				Parent:     office.Parent,
				Children:   append([]string(nil), office.Children...),
			},
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].score == items[j].score {
			return items[i].area.Code < items[j].area.Code
		}

		return items[i].score > items[j].score
	})

	responseItems := make([]gen.Area, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, item.area)
	}

	return gen.AreasResponse{Items: responseItems}
}

func matchAreaName(mode gen.AreaMatchMode, name string, query *string) (int, bool) {
	if query == nil {
		return 0, true
	}

	normalizedName := strings.TrimSpace(name)
	normalizedQuery := strings.TrimSpace(*query)
	if normalizedQuery == "" {
		return 0, true
	}

	switch mode {
	case "", gen.AreaMatchMode("exact"):
		return 100, normalizedName == normalizedQuery
	case gen.AreaMatchMode("prefix"):
		return 80, strings.HasPrefix(normalizedName, normalizedQuery)
	case gen.AreaMatchMode("partial"):
		return 60, strings.Contains(normalizedName, normalizedQuery)
	case gen.AreaMatchMode("suggested"):
		switch {
		case normalizedName == normalizedQuery:
			return 100, true
		case strings.HasPrefix(normalizedName, normalizedQuery):
			return 80, true
		case strings.Contains(normalizedName, normalizedQuery):
			return 60, true
		default:
			return 0, false
		}
	default:
		return 0, false
	}
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
