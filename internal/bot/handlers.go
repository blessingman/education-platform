package bot

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	model "github.com/blessingman/education-platform/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Обработчик команды /start
func handleStart(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil {
		welcomeText := "👋 Добро пожаловать! Пожалуйста, зарегистрируйтесь, отправив свой пропуск с помощью команды /register."
		msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
		b.Telegram.Send(msg)
		return
	}

	welcomeText := fmt.Sprintf("👋 Добро пожаловать, %s!\nВыберите одну из опций ниже:", user.Name)
	keyboard := getMainMenuKeyboard(user.Role)
	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ReplyMarkup = keyboard
	b.Telegram.Send(msg)
}

// Обработчик команды /register
func handleRegister(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err == nil {
		if user.Active {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Вы уже зарегистрированы и активны.")
			b.Telegram.Send(msg)
		} else {
			// Реактивируем аккаунт
			_, err = b.DB.Exec(`UPDATE users SET active = TRUE WHERE telegram_id = $1`, user.TelegramID)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при активации аккаунта. Попробуйте позже.")
				b.Telegram.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Ваш аккаунт был успешно активирован.")
				b.Telegram.Send(msg)
			}
		}
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, введите ваш пропуск для регистрации.")
	b.Telegram.Send(msg)
}

// Обработчик команды /schedule
func handleSchedule(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, сначала зарегистрируйтесь с помощью команды /register.")
		b.Telegram.Send(msg)
		return
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	var buttons [][]tgbotapi.InlineKeyboardButton

	if user.Role == "student" || user.Role == "teacher" || user.Role == "admin" {
		if user.Role == "teacher" || user.Role == "admin" {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📅 По группе", "schedule_group"),
				tgbotapi.NewInlineKeyboardButtonData("🔸 По преподавателю", "schedule_teacher"),
			))
		} else {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📅 По группе", "schedule_group"),
			))
		}

		keyboard = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите тип расписания:")
	msg.ReplyMarkup = keyboard
	b.Telegram.Send(msg)
}

// Обработчик команды /admin
func handleAdmin(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil || user.Role != "admin" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "У вас нет прав администратора.")
		b.Telegram.Send(msg)
		return
	}

	adminMenu := tgbotapi.NewMessage(message.Chat.ID, "Добро пожаловать в административное меню! Используйте следующие команды:\n"+
		"/add_passcode - Добавить пропуск\n"+
		"/add_student - Добавить студента\n"+
		"/add_teacher - Добавить преподавателя\n"+
		"/add_schedule - Добавить расписание")
	b.Telegram.Send(adminMenu)
}

// Обработчик команды /add_passcode
func handleAddPasscode(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil || user.Role != "admin" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "У вас нет прав администратора.")
		b.Telegram.Send(msg)
		return
	}

	args := strings.Split(message.CommandArguments(), " ")
	if len(args) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Используйте формат: /add_passcode <код> <роль>")
		b.Telegram.Send(msg)
		return
	}

	code := args[0]
	role := args[1]

	if role != "student" && role != "teacher" && role != "admin" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неверная роль. Используйте: student, teacher или admin.")
		b.Telegram.Send(msg)
		return
	}

	_, err = b.DB.Exec(`INSERT INTO passcodes (code, role) VALUES ($1, $2)`, code, role)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при добавлении пропуска.")
		b.Telegram.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Пропуск %s для роли %s успешно добавлен.", code, role))
	b.Telegram.Send(msg)
}

// Обработчик текстовых сообщений
func handleMessage(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil {
		handlePasscodeRegistration(b, message)
		return
	}

	msgText := fmt.Sprintf("Привет, %s! Я не понимаю эту команду. Пожалуйста, используйте доступные команды.", user.Name)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	b.Telegram.Send(msg)
}

// Функция для получения пользователя по Telegram ID
func getUserByTelegramID(db *sql.DB, telegramID int64) (*model.User, error) {
	var user model.User
	query := `SELECT id, telegram_id, name, email, role, active FROM users WHERE telegram_id = $1`
	err := db.QueryRow(query, telegramID).Scan(&user.ID, &user.TelegramID, &user.Name, &user.Email, &user.Role, &user.Active)
	if err != nil {
		return nil, err
	}

	// Проверяем, активен ли пользователь
	if !user.Active {
		return nil, fmt.Errorf("пользователь неактивен")
	}

	return &user, nil
}

// Функция для получения главного меню
func getMainMenuKeyboard(role string) tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton

	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📅 Просмотреть расписание"),
		tgbotapi.NewKeyboardButton("📚 Учебные материалы"),
	))

	if role == "admin" {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚙ Администрирование"),
		))
	}

	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📝 Оставить отзыв"),
	))

	return tgbotapi.NewReplyKeyboard(buttons...)
}

// Обработка регистрации через пропуск
func handlePasscodeRegistration(b *Bot, message *tgbotapi.Message) {
	passcode := strings.TrimSpace(message.Text)

	matched, err := regexp.MatchString(`^(ST|TE|AD)-[A-Z0-9]{5}$`, passcode)
	if err != nil || !matched {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Неверный формат пропуска. Пожалуйста, проверьте ваш пропуск и попробуйте снова.")
		b.Telegram.Send(msg)
		return
	}

	var role string
	var isUsed bool
	query := `SELECT role, is_used FROM passcodes WHERE code = $1`
	err = b.DB.QueryRow(query, passcode).Scan(&role, &isUsed)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Пропуск не найден. Пожалуйста, проверьте ваш пропуск или обратитесь к администратору.")
		b.Telegram.Send(msg)
		return
	}

	if isUsed {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Этот пропуск уже использован. Пожалуйста, получите новый пропуск от администратора.")
		b.Telegram.Send(msg)
		return
	}

	var assignedRole string
	switch {
	case strings.HasPrefix(passcode, "ST-"):
		assignedRole = "student"
	case strings.HasPrefix(passcode, "TE-"):
		assignedRole = "teacher"
	case strings.HasPrefix(passcode, "AD-"):
		assignedRole = "admin"
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Неверный тип пропуска.")
		b.Telegram.Send(msg)
		return
	}

	_, err = b.DB.Exec(`INSERT INTO users (telegram_id, name, email, role) VALUES ($1, $2, $3, $4)`,
		int64(message.From.ID), message.From.UserName, "", assignedRole)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Произошла ошибка при регистрации. Пожалуйста, попробуйте позже.")
		b.Telegram.Send(msg)
		log.Printf("Ошибка регистрации пользователя: %v", err)
		return
	}

	_, err = b.DB.Exec(`UPDATE passcodes SET is_used = TRUE, used_at = $1 WHERE code = $2`, time.Now(), passcode)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Произошла ошибка при обновлении пропуска. Пожалуйста, обратитесь к администратору.")
		b.Telegram.Send(msg)
		log.Printf("Ошибка обновления пропуска: %v", err)
		return
	}

	welcomeText := fmt.Sprintf("✅ Вы успешно зарегистрированы как *%s*. Добро пожаловать!", assignedRole)
	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ParseMode = "Markdown"
	b.Telegram.Send(msg)
}

// Обработка inline-кнопок
// Обработка inline-кнопок
func handleCallbackQuery(b *Bot, callback *tgbotapi.CallbackQuery) {
	var responseText string

	switch callback.Data {
	case "schedule_group":
		responseText = "Вы выбрали просмотр расписания по группе."
	case "schedule_teacher":
		responseText = "Вы выбрали просмотр расписания по преподавателю."
	default:
		responseText = "Неизвестный выбор."
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, responseText)
	b.Telegram.Send(msg)

	// Подтверждаем нажатие кнопки
	callbackResponse := tgbotapi.NewCallback(callback.ID, "Выбор принят")
	b.Telegram.AnswerCallbackQuery(callbackResponse)
}

// Обработка команды /logout
func handleLogout(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Вы не зарегистрированы.")
		b.Telegram.Send(msg)
		return
	}

	// Обновляем статус пользователя на "неактивен"
	_, err = b.DB.Exec(`UPDATE users SET active = FALSE WHERE telegram_id = $1`, user.TelegramID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при выходе из системы. Попробуйте позже.")
		b.Telegram.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы успешно вышли из системы. Чтобы снова войти, используйте команду /register и введите ваш пропуск.")
	b.Telegram.Send(msg)
}
