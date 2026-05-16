# Dockerfile
FROM golang:1.25.6-alpine AS builder

# Устанавливаем переменные окружения для оптимизации сборки
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости (кэшируется как отдельный слой)
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем приложение
RUN go build -ldflags="-s -w" -o auth-service ./cmd/auth/

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations
RUN chmod +x auth-service
CMD ["./auth-service"]