# Scout Education Service — Makefile
.PHONY: build test test-validator test-activity run clean

# Сборка бинарника
build:
	go build -o education-service ./cmd/education

# Запуск всех тестов
test:
	go test ./... -v

# Запуск только тестов валидации (курсы + активности)
test-validator:
	go test ./internal/validator/... -v

# Запуск тестов для курсов
test-course:
	go test ./internal/validator/... -run TestCreateCourse -v

# Запуск тестов для активностей
test-activity:
	go test ./internal/validator/... -run TestCreateActivity -v

# Запуск сервиса локально (без Docker)
run:
	go run ./cmd/education

# Очистка (удаление бинарника)
clean:
	rm -f education-service