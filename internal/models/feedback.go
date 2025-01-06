package models

import "time"

type Feedback struct {
	ID      int       `db:"id"`
	UserID  int       `db:"user_id"`
	Message string    `db:"message"`
	Created time.Time `db:"created_at"`
}
