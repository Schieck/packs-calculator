# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24.3 AS builder

RUN apt-get update && apt-get install -y git ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/server/main.go -o ./docs

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app cmd/server/main.go

# Final stage - using distroless for better security (non-root user)
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 8080

# Health check omitted - should be handled by orchestrator (K8s, Docker Compose, etc.)

CMD ["./app"] 