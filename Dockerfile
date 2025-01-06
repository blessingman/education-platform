# Используем официальный образ Go
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы проекта
COPY . .

# Устанавливаем зависимости
RUN go mod tidy

# Компилируем приложение
RUN go build -o bot ./cmd/bot/main.go

# Копируем SQL-скрипт для миграции базы данных
COPY scripts/setup_db.sql /docker-entrypoint-initdb.d/

# Запускаем приложение
CMD ["./bot"]
