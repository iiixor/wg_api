# Этап сборки
FROM golang:1.24-alpine AS builder

# Установка необходимых утилит
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем документацию Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go -d ./ --parseDependency --output ./docs
RUN apk add --no-cache \
    docker-cli \
    docker-compose

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o wg-api ./cmd/api

# Этап финальный
FROM alpine:latest

# Установка WireGuard и необходимых утилит
RUN apk add --no-cache wireguard-tools iptables ip6tables

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/wg-api .
COPY --from=builder /app/docs ./docs

# Порт для API
EXPOSE 8080

# Порт для WireGuard
EXPOSE 51820/udp

# Запуск приложения
CMD ["./wg-api"]
