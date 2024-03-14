package service

import (
	"github.com/artem-shestakov/telegram-budget/internal/repository"
)

type BudgetService struct {
	repository *repository.Repository
}

func NewBudgetService(repository *repository.Repository) *BudgetService {
	return &BudgetService{
		repository: repository,
	}
}

func (s *BudgetService) CreateBudget(chatId int64, title string) (int, error) {
	return s.repository.Budget.CreateBudget(chatId, title)
}
