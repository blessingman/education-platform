package models

import "time"

type Material struct {
	ID       int       `db:"id"`
	GroupID  int       `db:"group_id"`
	Title    string    `db:"title"`
	FileURL  string    `db:"file_url"`
	Uploaded time.Time `db:"uploaded_at"`
}
