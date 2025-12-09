# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o users-service .

# Dev stage with hot reload (Air)
FROM golang:1.25.4-alpine AS dev
WORKDIR /app
# Air module path has moved to github.com/air-verse/air
RUN apk add --no-cache git && go install github.com/air-verse/air@latest
COPY . .
CMD ["air", "-c", ".air.toml"]

# Prod stage
FROM alpine:latest AS prod

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

COPY --from=builder /app/users-service .

EXPOSE 5000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:5000/health || exit 1

# Run the application
CMD ["./users-service"]
