package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
}

func LoadConfig() (*Config, error) {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env файл не найден, используем переменные окружения")
	}

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
	}, nil
}
