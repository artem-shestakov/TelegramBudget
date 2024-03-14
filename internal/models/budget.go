package models

type Budget struct {
	Id    int    `db:"id"`
	Title string `db:"title"`
}
