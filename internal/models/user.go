package models

type User struct {
	ID         int    `db:"id"`
	TelegramID int64  `db:"telegram_id"`
	Name       string `db:"name"`
	Email      string `db:"email"`
	Role       string `db:"role"`
}
