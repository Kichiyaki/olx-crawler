package usecase

import (
	"olx-crawler/models"
	"olx-crawler/suggestion"
)

type usecase struct {
	suggestionRepo suggestion.Repository
}

func NewSuggestionUsecase(repo suggestion.Repository) suggestion.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Fetch(f *models.SuggestionFilter) (models.PaginatedResponse, error) {
	if f == nil {
		f = &models.SuggestionFilter{}
		f.Limit = 100
	}
	return ucase.suggestionRepo.Fetch(f)
}

func (ucase *usecase) GetByID(id uint) (*models.Suggestion, error) {
	return ucase.suggestionRepo.GetByID(id)
}

func (ucase *usecase) Delete(ids ...uint) ([]*models.Suggestion, error) {
	f := &models.SuggestionFilter{
		ID: ids,
	}
	pagination, err := ucase.suggestionRepo.Fetch(f)
	if err != nil {
		return nil, err
	}
	items := pagination.Items.([]*models.Suggestion)
	if err := ucase.suggestionRepo.Delete(f); err != nil {
		return nil, err
	}
	return items, nil
}
