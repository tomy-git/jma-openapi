<!--
SPDX-FileCopyrightText: 2026 The jma-openapi contributors

SPDX-License-Identifier: MPL-2.0
-->

# jma-openapi

Unofficial OpenAPI wrapper for the Japan Meteorological Agency (JMA) BOSAI JSON data.

For a Japanese reference translation, see [README-JA.md](README-JA.md).

## Overview

`jma-openapi` is an unofficial API wrapper project that provides a simpler and more developer-friendly interface for the Japan Meteorological Agency (JMA) BOSAI JSON endpoints.

The project normalizes JMA weather and disaster-related data and exposes it through a consistent REST API defined by an OpenAPI specification. This makes it easier for application developers to consume JMA data without dealing with the complexity of the original JSON structures.

This project is not affiliated with, endorsed by, or maintained by the Japan Meteorological Agency.

---

## Goals

- Provide a simple wrapper for JMA BOSAI JSON endpoints
- Offer a stable and developer-friendly REST API interface
- Publish an OpenAPI specification for implemented endpoints
- Reduce the complexity of consuming raw JMA JSON structures
- Serve as a foundation for applications that use Japanese weather and disaster data

---

## Architecture

The project follows a **spec-first API design** approach.

OpenAPI specifications define the API contract first, and server code is generated from the specification using Go tooling.

API handlers then implement the generated interfaces.

---

## Technology Stack

- Language: **Go**
- Router: **chi**
- OpenAPI: **oapi-codegen**
- HTTP layer: **net/http (Go standard library)**
- Testing: **go test**
- Lint: **golangci-lint**
- API docs UI: **bundled**
- Deployment: **Cloud Run**
- Logging: **log/slog**

This combination provides:

- High performance
- Minimal framework dependency
- Strong compatibility with Go ecosystem
- Long-term maintainability

---

## Development Policy

- Follow **GitHub Flow** for development
- Keep `main` always in a working state
- Use **OpenAPI specification as the source of truth**
- Generate server interfaces using `oapi-codegen`
- Generate types, server interfaces, and embedded spec bindings into dedicated generated files
- Implement business logic in handlers separate from generated code
- Bundle an OpenAPI UI for human-friendly API reference
- Run tests with `go test`
- Run lint with `golangci-lint`
- Target Cloud Run as the initial deployment platform
- Use `log/slog` for structured logging
- Build and deploy the service with a `Dockerfile`
- Configure Cloud Run minimum instances as `0` for the initial release
- Keep the architecture simple and modular

---

## Scope

The initial scope of the project includes the following:

- Health check endpoint
- Forecast endpoint
- Area metadata endpoint
- Basic OpenAPI specification

Additional endpoints may be added later for:

- Overview forecasts
- Weather warnings
- Earthquake information
- Tsunami alerts
- Other JMA BOSAI data resources

---

## Getting Started

### Prerequisites

- Go `1.26.1`
- Docker
- `oapi-codegen`
- `golangci-lint`

If you use `mise`, you can align the Go version with the repository's `mise.toml`.

### Local Run

Install dependencies:

```bash
go mod tidy
```

Run tests and lint:

```bash
go test ./...
golangci-lint run
```

GitHub Actions runs `go test ./...`, `golangci-lint`, REUSE license checks, `govulncheck ./...`, and `go build ./...` on pull requests and pushes to `main`.

Start the server:

```bash
go run ./cmd/server
```

Endpoints to check:

- `http://localhost:8080/healthz`
- `http://localhost:8080/v1/areas`
- `http://localhost:8080/v1/areas/130000`
- `http://localhost:8080/v1/forecasts/130000`
- `http://localhost:8080/openapi.yaml`
- `http://localhost:8080/openapi.json`
- `http://localhost:8080/docs`

### Environment Variables

- `PORT`
  - Default: `8080`
- `SERVICE_NAME`
  - Default: `jma-openapi`
- `APP_VERSION`
  - Default: `dev`
- `JMA_BASE_URL`
  - Default: `https://www.jma.go.jp/bosai`
- `LOG_FORMAT`
  - Default: `text`
  - Use `json` for structured logs

### Regenerate OpenAPI Bindings

```bash
oapi-codegen -config openapi/oapi-codegen-types.yaml openapi/openapi.yaml
oapi-codegen -config openapi/oapi-codegen-server.yaml openapi/openapi.yaml
oapi-codegen -config openapi/oapi-codegen-spec.yaml openapi/openapi.yaml
```

### License Headers

Add or refresh SPDX headers and `.license` sidecar files:

```bash
mise run license-annotate
```

Verify REUSE compliance:

```bash
mise run license-lint
```

### Run with Docker

Start the service with Docker Compose:

```bash
docker compose up --build
```

The app is exposed at `http://localhost:8080`.

`docker compose build` and `docker compose up --build` run `go test ./...` and `golangci-lint run` during the image build. The container starts only when both checks pass.

### Cloud Run

Build and push the image:

```bash
docker build -t gcr.io/<PROJECT_ID>/jma-openapi:latest .
docker push gcr.io/<PROJECT_ID>/jma-openapi:latest
```

Apply the service definition:

```bash
gcloud run services replace deploy/cloudrun/service.yaml --region <REGION>
```

---

## Disclaimer

This is an unofficial project. The upstream JMA BOSAI JSON format and endpoint behavior may change at any time. Compatibility is not guaranteed.

---

## License

This project is licensed under the Mozilla Public License 2.0 (MPL-2.0).
