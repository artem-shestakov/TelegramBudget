package models

type Budget struct {
	id    int    `db:"id"`
	title string `db:"title"`
}
