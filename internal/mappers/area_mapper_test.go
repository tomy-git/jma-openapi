// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package mappers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/tomy-git/jma-openapi/internal/clients"
	"github.com/tomy-git/jma-openapi/internal/gen"
)

func TestAreaMapper_ToAreasResponse(t *testing.T) {
	t.Parallel()

	var document clients.AreaJSONDocument

	payload, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "area.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(payload, &document); err != nil {
		t.Fatal(err)
	}

	mapper := NewAreaMapper()
	tests := []struct {
		name       string
		filter     AreaFilter
		wantCount  int
		wantFirst  string
		wantParent string
	}{
		{
			name: "parent",
			filter: AreaFilter{
				Parent: stringPtr("010300"),
			},
			wantCount:  9,
			wantFirst:  "080000",
			wantParent: "010300",
		},
		{
			name: "name",
			filter: AreaFilter{
				Name:     stringPtr("東京都"),
				NameMode: gen.Exact,
			},
			wantCount: 1,
			wantFirst: "130000",
		},
		{
			name: "name prefix",
			filter: AreaFilter{
				Name:     stringPtr("東京"),
				NameMode: gen.Prefix,
			},
			wantCount: 1,
			wantFirst: "130000",
		},
		{
			name: "name partial",
			filter: AreaFilter{
				Name:     stringPtr("京都"),
				NameMode: gen.Partial,
			},
			wantCount: 2,
			wantFirst: "130000",
		},
		{
			name: "name suggested",
			filter: AreaFilter{
				Name:     stringPtr("京"),
				NameMode: gen.Suggested,
			},
			wantCount: 2,
			wantFirst: "260000",
		},
		{
			name: "officeName",
			filter: AreaFilter{
				OfficeName: stringPtr("気象庁"),
			},
			wantCount: 1,
			wantFirst: "130000",
		},
		{
			name: "child",
			filter: AreaFilter{
				Child: stringPtr("130010"),
			},
			wantCount: 1,
			wantFirst: "130000",
		},
		{
			name: "combined",
			filter: AreaFilter{
				Parent:     stringPtr("010300"),
				Name:       stringPtr("東京都"),
				OfficeName: stringPtr("気象庁"),
				Child:      stringPtr("130010"),
			},
			wantCount: 1,
			wantFirst: "130000",
		},
		{
			name: "empty",
			filter: AreaFilter{
				Name: stringPtr("存在しない地域"),
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			response := mapper.ToAreasResponse(document, tt.filter)
			if len(response.Items) != tt.wantCount {
				t.Fatalf("expected %d items, got %d", tt.wantCount, len(response.Items))
			}
			if tt.wantCount == 0 {
				return
			}
			if response.Items[0].Code != tt.wantFirst {
				t.Fatalf("expected first code %s, got %s", tt.wantFirst, response.Items[0].Code)
			}
			if tt.wantParent != "" && response.Items[0].Parent != tt.wantParent {
				t.Fatalf("expected parent %s, got %s", tt.wantParent, response.Items[0].Parent)
			}
		})
	}
}

func stringPtr(value string) *string {
	return &value
}
