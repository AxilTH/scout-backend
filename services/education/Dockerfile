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
RUN go build -ldflags="-s -w" -o education-service ./cmd/education

# Скачиваем wait-for-it.sh
RUN wget -O wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
    && chmod +x wait-for-it.sh

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=builder /app/education-service .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/wait-for-it.sh .
CMD ["./wait-for-it.sh", "postgresql:5432", "--", "./education-service"]