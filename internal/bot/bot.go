package bot

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot - структура для хранения бота и базы данных
type Bot struct {
	Telegram      *tgbotapi.BotAPI
	DB            *sql.DB
	UpdateChannel tgbotapi.UpdatesChannel // Добавляем канал для получения обновлений
}

// NewBot - функция для создания нового бота
func NewBot(token string, db *sql.DB) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// Создаем канал для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := botAPI.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Telegram:      botAPI,
		DB:            db,
		UpdateChannel: updates, // Инициализируем канал обновлений
	}, nil
}

// Start запускает бота
// Start запускает бота
func (b *Bot) Start() {
	for update := range b.UpdateChannel {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(b, update.Message)
			} else {
				handleMessage(b, update.Message)
			}
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(b, update.CallbackQuery)
		}
	}
}

// Обработчик всех команд
func handleCommand(b *Bot, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		handleStart(b, message)
	case "register":
		handleRegister(b, message)
	case "schedule":
		handleSchedule(b, message)
	case "admin":
		handleAdmin(b, message)
	case "add_passcode":
		handleAddPasscode(b, message)
	case "logout":
		handleLogout(b, message)
	case "add_schedule":
		handleAddSchedule(b, message)

	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Используйте /start или /help.")
		b.Telegram.Send(msg)
	}
}
