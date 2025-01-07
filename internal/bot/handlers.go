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
		// Пользователь не найден, предлагаем зарегистрироваться
		welcomeText := "👋 Добро пожаловать! Пожалуйста, зарегистрируйтесь, отправив свой пропуск с помощью команды /register."
		msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
		b.Telegram.Send(msg)
		return
	}

	// Пользователь найден и активен
	if user.Active {
		welcomeText := fmt.Sprintf("👋 Добро пожаловать обратно, %s!\nВыберите одну из опций ниже:", user.Name)
		keyboard := getMainMenuKeyboard(user.Role) // Меню будет зависеть от роли пользователя
		msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
		msg.ReplyMarkup = keyboard
		b.Telegram.Send(msg)
	} else {
		// Пользователь найден, но неактивен
		msg := tgbotapi.NewMessage(message.Chat.ID, "🔔 Ваш аккаунт был деактивирован. Пожалуйста, используйте команду /register для повторной активации.")
		b.Telegram.Send(msg)
	}
}

// Обработчик команды /register
func handleRegister(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err == nil {
		if user.Active {
			// Пользователь уже зарегистрирован и активен
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вы уже зарегистрированы как %s.", user.Role))
			b.Telegram.Send(msg)
		} else {
			// Реактивируем аккаунт
			_, err = b.DB.Exec(`UPDATE users SET active = TRUE WHERE telegram_id = $1`, user.TelegramID)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Ошибка при активации аккаунта. Попробуйте позже.")
				b.Telegram.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "✅ Ваш аккаунт успешно активирован.")
				b.Telegram.Send(msg)
			}
		}
		return
	}

	// Если пользователь не найден, просим ввести пропуск
	msg := tgbotapi.NewMessage(message.Chat.ID, "🔑 Пожалуйста, введите ваш пропуск для регистрации.")
	b.Telegram.Send(msg)
}

// Обработчик команды /admin
func handleAdmin(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Вы не зарегистрированы. Используйте /register для регистрации.")
		b.Telegram.Send(msg)
		return
	}

	if user.Role != "admin" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ У вас нет прав администратора.")
		b.Telegram.Send(msg)
		return
	}

	adminMenu := tgbotapi.NewMessage(message.Chat.ID, "⚙️ Добро пожаловать в административное меню! Используйте следующие команды:\n"+
		"/add_passcode - Добавить пропуск\n"+
		"/add_schedule - Добавить расписание\n"+
		"/logout - Выйти из системы")
	b.Telegram.Send(adminMenu)
}

// Обработчик текстовых сообщений
func handleMessage(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil {
		handlePasscodeRegistration(b, message)
		return
	}

	switch message.Text {
	case "⚙ Администрирование":
		if user.Role == "admin" {
			handleAdmin(b, message)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "❌ У вас нет прав администратора.")
			b.Telegram.Send(msg)
		}
	case "📅 Просмотреть расписание":
		handleSchedule(b, message)
	case "📚 Список предметов":
		handleListSubjects(b, message) // Новая команда для вывода списка предметов
	case "👥 Список групп":
		handleListGroups(b, message) // Новая команда для вывода списка групп
	default:
		msgText := fmt.Sprintf("Привет, %s! Я не понимаю эту команду. Используйте доступные команды:\n"+
			"- 📅 Просмотреть расписание\n"+
			"- 📚 Список предметов\n"+
			"- 👥 Список групп", user.Name)
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		b.Telegram.Send(msg)
	}
}

// Функция для получения пользователя по Telegram ID
func getUserByTelegramID(db *sql.DB, telegramID int64) (*model.User, error) {
	var user model.User
	query := `SELECT id, telegram_id, name, email, role, active FROM users WHERE telegram_id = $1`
	err := db.QueryRow(query, telegramID).Scan(&user.ID, &user.TelegramID, &user.Name, &user.Email, &user.Role, &user.Active)
	if err != nil {
		return nil, err
	}

	// Пользователь найден, возвращаем его
	return &user, nil
}

// Функция для получения главного меню
func getMainMenuKeyboard(role string) tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton

	// Общие кнопки для всех пользователей
	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📅 Просмотреть расписание"),
	))

	// Кнопка "Администрирование" доступна только администраторам
	if role == "admin" {
		adminButtons := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚙ Администрирование"),
		)
		buttons = append(buttons, adminButtons)
	}

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
		// Здесь можно добавить логику для вывода расписания по группе
	case "schedule_teacher":
		responseText = "Вы выбрали просмотр расписания по преподавателю."
		// Здесь можно добавить логику для вывода расписания по преподавателю
	case "list_subjects":
		responseText = "Список предметов:\n"
		subjects, err := getSubjects(b.DB)
		if err != nil {
			responseText = "❌ Ошибка при получении списка предметов."
		} else {
			for _, subject := range subjects {
				responseText += fmt.Sprintf("- %s\n", subject.Name)
			}
		}
	case "list_groups":
		responseText = "Список групп:\n"
		groups, err := getGroups(b.DB)
		if err != nil {
			responseText = "❌ Ошибка при получении списка групп."
		} else {
			for _, group := range groups {
				responseText += fmt.Sprintf("- %s\n", group.Name)
			}
		}
	default:
		responseText = "Неизвестный выбор."
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, responseText)
	b.Telegram.Send(msg)

	// Подтверждаем нажатие кнопки
	callbackResponse := tgbotapi.NewCallback(callback.ID, "Выбор принят")
	b.Telegram.AnswerCallbackQuery(callbackResponse)
}

func handleListSubjects(b *Bot, message *tgbotapi.Message) {
	subjects, err := getSubjects(b.DB)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Ошибка при получении списка предметов.")
		b.Telegram.Send(msg)
		return
	}

	var subjectList string
	for _, subject := range subjects {
		subjectList += fmt.Sprintf("- %s\n", subject.Name)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "📚 Список предметов:\n"+subjectList)
	b.Telegram.Send(msg)
}

func handleListGroups(b *Bot, message *tgbotapi.Message) {
	groups, err := getGroups(b.DB)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Ошибка при получении списка групп.")
		b.Telegram.Send(msg)
		return
	}

	var groupList string
	for _, group := range groups {
		groupList += fmt.Sprintf("- %s\n", group.Name)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "👥 Список групп:\n"+groupList)
	b.Telegram.Send(msg)
}

// Обработка команды /logout
func handleLogout(b *Bot, message *tgbotapi.Message) {
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Вы не зарегистрированы.")
		b.Telegram.Send(msg)
		return
	}

	// Обновляем время последнего выхода и статус
	_, err = b.DB.Exec(`UPDATE users SET active = FALSE, last_logout = $1 WHERE telegram_id = $2`, time.Now(), user.TelegramID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при выходе из системы. Попробуйте позже.")
		b.Telegram.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы успешно вышли из системы. Чтобы снова войти, используйте команду /start.")
	b.Telegram.Send(msg)
}

// Обработчик команды
// добавить подсказки
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

func handleAddSchedule(b *Bot, message *tgbotapi.Message) {
	// Проверяем, является ли пользователь администратором
	user, err := getUserByTelegramID(b.DB, int64(message.From.ID))
	if err != nil || user.Role != "admin" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ У вас нет прав администратора.")
		b.Telegram.Send(msg)
		return
	}

	// Шаг 1: Получаем список групп
	groups := make(map[string]int)
	rows, err := b.DB.Query("SELECT id, name FROM groups")
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Ошибка при получении списка групп.")
		b.Telegram.Send(msg)
		return
	}
	defer rows.Close()

	var groupOptions string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err == nil {
			groups[name] = id
			groupOptions += fmt.Sprintf("- %s\n", name)
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "📅 Выберите группу:\n"+groupOptions)
	b.Telegram.Send(msg)

	// Ждем ответ с названием группы
	var groupName string
	for update := range b.UpdateChannel {
		if update.Message != nil && update.Message.Chat.ID == message.Chat.ID {
			groupName = update.Message.Text
			break
		}
	}

	groupID, ok := groups[groupName]
	if !ok {
		b.Telegram.Send(tgbotapi.NewMessage(message.Chat.ID, "❌ Группа не найдена."))
		return
	}

	// Шаг 2: Получаем список предметов
	subjects := make(map[string]int)
	rows, err = b.DB.Query("SELECT id, name FROM subjects")
	if err != nil {
		b.Telegram.Send(tgbotapi.NewMessage(message.Chat.ID, "❌ Ошибка при получении списка предметов."))
		return
	}
	defer rows.Close()

	var subjectOptions string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err == nil {
			subjects[name] = id
			subjectOptions += fmt.Sprintf("- %s\n", name)
		}
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "📚 Выберите предмет:\n"+subjectOptions)
	b.Telegram.Send(msg)

	// Ждем ответ с названием предмета
	var subjectName string
	for update := range b.UpdateChannel {
		if update.Message != nil && update.Message.Chat.ID == message.Chat.ID {
			subjectName = update.Message.Text
			break
		}
	}

	subjectID, ok := subjects[subjectName]
	if !ok {
		b.Telegram.Send(tgbotapi.NewMessage(message.Chat.ID, "❌ Предмет не найден."))
		return
	}

	// Шаг 3: Запрашиваем дату и время
	b.Telegram.Send(tgbotapi.NewMessage(message.Chat.ID, "⏰ Введите дату и время в формате: YYYY-MM-DD HH:MM"))

	// Ждем ответ с датой и временем
	var dateTimeStr string
	for update := range b.UpdateChannel {
		if update.Message != nil && update.Message.Chat.ID == message.Chat.ID {
			dateTimeStr = update.Message.Text
			break
		}
	}

	// Парсим дату и время
	parsedDateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Неверный формат даты и времени. Используйте формат: YYYY-MM-DD HH:MM")
		b.Telegram.Send(msg)
		return
	}

	// Шаг 4: Проверяем уникальность расписания
	var existingScheduleID int
	err = b.DB.QueryRow("SELECT id FROM schedules WHERE group_id = $1 AND subject_id = $2 AND date = $3 AND time = $4", groupID, subjectID, parsedDateTime.Format("2006-01-02"), parsedDateTime.Format("15:04")).Scan(&existingScheduleID)
	if err == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Расписание для этой группы и предмета в указанное время уже существует.")
		b.Telegram.Send(msg)
		return
	}

	// Шаг 5: Добавляем расписание в базу данных
	_, err = b.DB.Exec("INSERT INTO schedules (group_id, subject_id, date, time) VALUES ($1, $2, $3, $4)", groupID, subjectID, parsedDateTime.Format("2006-01-02"), parsedDateTime.Format("15:04"))
	if err != nil {
		b.Telegram.Send(tgbotapi.NewMessage(message.Chat.ID, "❌ Ошибка при добавлении расписания."))
		return
	}

	b.Telegram.Send(tgbotapi.NewMessage(message.Chat.ID, "✅ Расписание успешно добавлено!"))
}

// getSubjects возвращает список всех предметов из базы данных
func getSubjects(db *sql.DB) ([]model.Subject, error) {
	var subjects []model.Subject

	rows, err := db.Query("SELECT id, name FROM subjects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subject model.Subject
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

// getGroups возвращает список всех групп из базы данных
func getGroups(db *sql.DB) ([]model.Group, error) {
	var groups []model.Group

	rows, err := db.Query("SELECT id, name FROM groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group model.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}
