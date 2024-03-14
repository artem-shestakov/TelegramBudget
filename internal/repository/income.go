package repository

import (
	"database/sql"
	"fmt"

	"github.com/artem-shestakov/telegram-budget/internal/models"
)

type IncomeRepository struct {
	db *sql.DB
}

func NewIncomeRepository(db *sql.DB) *IncomeRepository {
	return &IncomeRepository{
		db: db,
	}
}

func (r *IncomeRepository) Create(income models.Income) (int, error) {
	var incomeId int
	query := fmt.Sprintf("INSERT INTO %s (title, plan, budget_id) values ($1, $2, $3) RETURNING id", incomesTable)
	row := r.db.QueryRow(query, income.Title, income.Plan, income.BudgetId)
	if err := row.Scan(&incomeId); err != nil {
		return 0, err
	}
	return incomeId, nil
}

func (r *IncomeRepository) GetAll(chatId int64) ([]models.Income, error) {
	var incomes []models.Income

	query := fmt.Sprintf("SELECT * FROM %s WHERE budget_id=$1", incomesTable)
	rows, err := r.db.Query(query, chatId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var income models.Income
		if err := rows.Scan(&income.ID, &income.Title, &income.Plan, &income.BudgetId); err != nil {
			return incomes, err
		}
		incomes = append(incomes, income)
	}
	if err := rows.Err(); err != nil {
		return incomes, err
	}
	return incomes, nil
}
