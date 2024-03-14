package models

type Income struct {
	ID       int     `db:"id"`
	Title    string  `db:"title"`
	Plan     float32 `db:"plan"`
	BudgetId int64   `db:"bidget_id"`
}
