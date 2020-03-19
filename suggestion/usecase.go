package suggestion

import (
	"olx-crawler/models"
)

type Usecase interface {
	Fetch(f *models.SuggestionFilter) (models.PaginatedResponse, error)
	GetByID(id uint) (*models.Suggestion, error)
	Delete(ids ...uint) ([]*models.Suggestion, error)
}
