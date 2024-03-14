package service

import (
	"github.com/artem-shestakov/telegram-budget/internal/models"
	"github.com/artem-shestakov/telegram-budget/internal/repository"
)

type IncomeService struct {
	repository *repository.Repository
}

func NewIncomeService(repository *repository.Repository) *IncomeService {
	return &IncomeService{
		repository: repository,
	}
}

func (s *IncomeService) Create(income models.Income) (int, error) {
	return s.repository.Income.Create(income)
}

func (s *IncomeService) GetAll(chatId int64) ([]models.Income, error) {
	return s.repository.Income.GetAll(chatId)
}
