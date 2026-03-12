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
	parent := "010300"
	response := mapper.ToAreasResponse(document, &parent)

	if len(response.Items) == 0 {
		t.Fatal("expected filtered area items")
	}
	if response.Items[0].Parent != parent {
		t.Fatalf("expected parent %s, got %s", parent, response.Items[0].Parent)
	}
}
