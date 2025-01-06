package bot

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot структура бота
type Bot struct {
	Telegram *tgbotapi.BotAPI
	DB       *sql.DB
}

// NewBot инициализирует новый бот
func NewBot(token string, db *sql.DB) (*Bot, error) {
	telegramBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	telegramBot.Debug = true
	log.Printf("Авторизовались как %s", telegramBot.Self.UserName)

	return &Bot{
		Telegram: telegramBot,
		DB:       db,
	}, nil
}

// Start запускает бота и начинает прослушивание обновлений
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.Telegram.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Ошибка получения обновлений: %v", err)
	}

	for update := range updates {
		if update.Message == nil { // Игнорируем любые обновления, которые не являются сообщениями
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handleStart(b, update.Message)
			case "register":
				handleRegister(b, update.Message)
			case "schedule":
				handleSchedule(b, update.Message)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда. Используйте /start или /register.")
				b.Telegram.Send(msg)
			}
		} else {
			// Обработка сообщений, не являющихся командами
			handleMessage(b, update.Message)
		}
	}
}
