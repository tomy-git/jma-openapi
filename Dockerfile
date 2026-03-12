# SPDX-FileCopyrightText: 2026 The jma-openapi contributors
#
# SPDX-License-Identifier: MPL-2.0

FROM golang:1.26.1-bookworm AS base
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go env -w GOPATH=/go && go mod download
COPY . .

FROM base AS verify
RUN CGO_ENABLED=0 go test ./...
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3
RUN /go/bin/golangci-lint run

FROM verify AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /workspace/bin/jma-openapi ./cmd/server

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /workspace/bin/jma-openapi /usr/local/bin/jma-openapi
COPY --from=builder /workspace/openapi /app/openapi
COPY --from=builder /workspace/web /app/web
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/jma-openapi"]
