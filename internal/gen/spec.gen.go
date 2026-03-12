// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

// Code generated for the current OpenAPI contract scaffold. DO NOT EDIT.

package gen

import "os"

const (
	OpenAPISpecYAMLPath = "openapi/openapi.yaml"
	OpenAPISpecJSONPath = "openapi/openapi.json"
)

func LoadOpenAPISpecYAML() ([]byte, error) {
	return os.ReadFile(OpenAPISpecYAMLPath)
}

func LoadOpenAPISpecJSON() ([]byte, error) {
	return os.ReadFile(OpenAPISpecJSONPath)
}
