package models

import "time"

type Passcode struct {
	Code   string    `db:"code"`
	Role   string    `db:"role"` // Возможные значения: student, teacher, admin
	IsUsed bool      `db:"is_used"`
	UsedAt time.Time `db:"used_at"`
}
