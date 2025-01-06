package models

type Teacher struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
