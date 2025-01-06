package models

import "time"

type Announcement struct {
	ID      int       `db:"id"`
	Title   string    `db:"title"`
	Content string    `db:"content"`
	Created time.Time `db:"created_at"`
}
