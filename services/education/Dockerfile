# Dockerfile
FROM golang:1.25.6-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod init education-service || true
RUN go build -o education-service .

# Минимальный образ для запуска
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/education-service .
CMD ["./education-service"]