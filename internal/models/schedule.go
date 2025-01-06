package models

import "time"

type Schedule struct {
	ID        int       `db:"id"`
	GroupID   int       `db:"group_id"`
	TeacherID int       `db:"teacher_id"`
	Subject   string    `db:"subject"`
	Time      time.Time `db:"time"`
	Location  string    `db:"location"`
}
