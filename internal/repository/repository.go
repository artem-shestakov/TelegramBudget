package repository

import (
	"database/sql"

	"github.com/artem-shestakov/telegram-budget/internal/models"
)

const (
	budgetTable  = "budgets"
	incomesTable = "incomes"
	topUpTable   = "top_ups"
)

type Budget interface {
	CreateBudget(chatId int64, title string) (int, error)
}

type Income interface {
	Create(income models.Income) (int, error)
	GetAll(chatId int64) ([]models.Income, error)
	TopUp(topUp models.TopUp) (int, error)
}

type Repository struct {
	Budget
	Income
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Budget: NewBudgetRepository(db),
		Income: NewIncomeRepository(db),
	}
}
