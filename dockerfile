FROM golang:alpine as builder

WORKDIR /app

COPY src .

RUN go mod init main

RUN go mod tidy

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder app .

ENTRYPOINT ./app -host=${DB_HOST} -port=${DB_PORT} -user=${DB_USER} -password=${DB_PASSWORD} -dbname=${DB_NAME}