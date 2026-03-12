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
