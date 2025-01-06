package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DatabaseURL   string
}

func LoadConfig() (*Config, error) {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env файл не найден, используем переменные окружения")
	}

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DatabaseURL:   getDatabaseURL(),
	}, nil
}

func getDatabaseURL() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=disable"
}
