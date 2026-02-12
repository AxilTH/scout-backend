# Dockerfile
FROM golang:1.25.6-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o education-service ./cmd/education

# Скачиваем wait-for-it.sh прямо в билд-стадии
RUN wget -O wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
    && chmod +x wait-for-it.sh

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=builder /app/education-service .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/wait-for-it.sh .
CMD ["./wait-for-it.sh", "postgresql:5432", "--", "./education-service"]