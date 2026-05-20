# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS go-builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/app

# Stage 2: Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY --from=go-builder /usr/local/bin/app /usr/local/bin/app
COPY --from=go-builder /app/internal/app/infrastructure/database/migrations ./internal/app/infrastructure/database/migrations

EXPOSE 8080

CMD ["/usr/local/bin/app"]
