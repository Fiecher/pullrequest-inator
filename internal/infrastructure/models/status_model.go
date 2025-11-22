package models

type Status struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
