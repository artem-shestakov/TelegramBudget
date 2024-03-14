package repository

import (
	"database/sql"
)

const (
	budgetTable = "budget"
)

type Budget interface {
	CreateBudget(chatId int64, title string) (int, error)
}

type Repository struct {
	Budget
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Budget: NewBudgetRepository(db),
	}
}
