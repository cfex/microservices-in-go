FROM golang:1.25.5-alpne as builder

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linus go build -a -o /app/services/notification-service/tmp/notification-service ./services/notfication-service

FROM alpine:3.19 as runtine
WORKDIR /app/services/notification-service
COPY --from=builder /app/services/notification-service/tmp/notfication-service ./tmp/notification-service
EXPOSE 5200
ENTRYPOINT [ "./tmp/notification-service" ]

FROM golang:1.25.5-alpine as dev 
WORKDIR /app
RUN apk add --no-cache git && go install github.com/air-verse/air@latest
COPY . .
EXPOSE 5200
CMD [ "air", "-c", "./services/notification-service/.air.toml" ]
