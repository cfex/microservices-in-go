# Build stage
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/services/projects-service/tmp/projects-service ./services/projects-service

FROM alpine:3.19 AS runtime
WORKDIR /app/services/projects-service
COPY --from=builder /app/services/projects-service/tmp/projects-service ./tmp/projects-service
EXPOSE 5000 50051
ENTRYPOINT ["./tmp/projects-service"]

FROM golang:1.25.5-alpine AS dev
WORKDIR /app
ENV TZ=Europe/Belgrade
RUN apk add --no-cache git && go install github.com/air-verse/air@latest
COPY . .
EXPOSE 5000 50051
CMD ["air", "-c", "./services/projects-service/.air.toml"]
