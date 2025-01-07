package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int          `json:"id"`
	TelegramID int64        `json:"telegram_id"`
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Role       string       `json:"role"`
	Active     bool         `json:"active"`
	LastLogout sql.NullTime `json:"last_logout"`
}

type Passcode struct {
	Code   string    `db:"code"`
	Role   string    `db:"role"` // Возможные значения: student, teacher, admin
	IsUsed bool      `db:"is_used"`
	UsedAt time.Time `db:"used_at"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Schedule struct {
	ID        int    `json:"id"`
	GroupID   int    `json:"group_id"`
	SubjectID int    `json:"subject_id"`
	Date      string `json:"date"`
	Time      string `json:"time"`
}

type Subject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Teacher struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
