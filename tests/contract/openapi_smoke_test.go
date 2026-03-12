// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package contract

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenAPIFilesExist(t *testing.T) {
	t.Parallel()

	for _, path := range []string{
		filepath.Join("..", "..", "openapi", "openapi.yaml"),
		filepath.Join("..", "..", "openapi", "openapi.json"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}
