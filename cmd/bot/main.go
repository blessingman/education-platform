package main

import (
	"log"

	"github.com/blessingman/education-platform/internal/bot"
	"github.com/blessingman/education-platform/internal/config"
	"github.com/blessingman/education-platform/internal/database"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к базе данных
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация и запуск бота
	telegramBot, err := bot.NewBot(cfg.TelegramToken, db)
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}

	telegramBot.Start()
}
