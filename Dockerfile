# Используем официальный образ Go для сборки приложения
FROM golang:latest AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Устанавливаем зависимости и компилируем приложение
RUN go mod tidy
RUN go build -o bot ./cmd/bot/main.go

# Используем образ Ubuntu для запуска приложения
FROM ubuntu:22.04

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y \
    libc6 \
    ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/bot .


# Копируем .env файл в контейнер
COPY .env /app/.env

# Копируем SQL-скрипт для базы данных
COPY ./scripts/setup_db.sql /docker-entrypoint-initdb.d/

# Запускаем приложение
CMD ["./bot"]
