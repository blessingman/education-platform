package models

type Subject struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
