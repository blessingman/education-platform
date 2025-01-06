package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/blessingman/education-platform/internal/bot"
	"github.com/blessingman/education-platform/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	))
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
