module github.com/tomy-git/jma-openapi

go 1.26.0

toolchain go1.26.1

require (
	github.com/go-chi/chi/v5 v5.2.5
	github.com/golangci/golangci-lint/v2 v2.11.3
	github.com/oapi-codegen/oapi-codegen/v2 v2.6.0
)

tool (
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint
	github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
)
