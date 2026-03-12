FROM golang:1.24-bullseye AS builder
WORKDIR /workspace
COPY go.mod ./
RUN go env -w GOPATH=/go && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /workspace/bin/jma-openapi ./cmd/server

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /workspace/bin/jma-openapi /usr/local/bin/jma-openapi
COPY --from=builder /workspace/openapi /app/openapi
COPY --from=builder /workspace/web /app/web
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/jma-openapi"]
