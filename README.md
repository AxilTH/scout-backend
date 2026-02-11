# Scout Backend

Микросервисная платформа для учёта профессионального пути бойцов студенческих педагогических отрядов (СПО).

## 🧱 Архитектура

- **Микросервисы**: Go (Gin)
- **База данных**: PostgreSQL 17
- **Файловое хранилище**: MinIO (S3-совместимое)
- **API Gateway**: Nginx
- **Контейнеризация**: Docker + Docker Compose

## 🚀 Быстрый старт

1. Установи [Docker](https://www.docker.com/) и [Docker Compose](https://docs.docker.com/compose/).
2. Склонируй репозиторий:
   ```bash
   git clone https://github.com/ваш-логин/scout-backend.git
   cd scout-backend
   ```
3. Создай .env из шаблона и заполни его:
   ```bash
   cp .env.example .env
   ```
4. Запусти проект:
   ```bash
   docker-compose up --build
   ```