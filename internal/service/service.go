package service

import (
	"github.com/artem-shestakov/telegram-budget/internal/repository"
)

type Budget interface {
	CreateBudget(chatId int64, title string) (int, error)
}

type Service struct {
	Budget
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Budget: NewBudgetService(repository),
	}
}
