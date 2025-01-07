package models

type User struct {
	ID         int    `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Active     bool   `json:"active"`
}
