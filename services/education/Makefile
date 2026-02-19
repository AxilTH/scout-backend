# Scout Education Service — Makefile

.PHONY: build test test-validator run clean

# Сборка бинарника
build:
	go build -o education-service ./cmd/education

# Запуск всех тестов
test:
	go test ./... -v

# Запуск только тестов валидации
test-validator:
	go test ./internal/validator/... -v

# Запуск сервиса локально (без Docker)
run:
	go run ./cmd/education

# Очистка (удаление бинарника)
clean:
	rm -f education-service
