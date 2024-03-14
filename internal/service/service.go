package service

import (
	"github.com/artem-shestakov/telegram-budget/internal/models"
	"github.com/artem-shestakov/telegram-budget/internal/repository"
)

type Budget interface {
	CreateBudget(chatId int64, title string) (int, error)
}

type Income interface {
	Create(income models.Income) (int, error)
	GetAll(chatId int64) ([]models.Income, error)
}

type Service struct {
	Budget
	Income
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Budget: NewBudgetService(repository),
		Income: NewIncomeService(repository),
	}
}
