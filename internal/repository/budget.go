package repository

import (
	"database/sql"
	"fmt"
)

type BudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{
		db: db,
	}
}

func (r *BudgetRepository) CreateBudget(chatId int64, title string) (int, error) {
	var budgetId int
	query := fmt.Sprintf("INSERT INTO %s (id, title) values ($1, $2) RETURNING id", budgetTable)
	row := r.db.QueryRow(query, chatId, title)
	if err := row.Scan(&budgetId); err != nil {
		return 0, err
	}
	return budgetId, nil
}
