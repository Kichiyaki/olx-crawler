package suggestion

import (
	"olx-crawler/models"
)

type Repository interface {
	Fetch(f *models.SuggestionFilter) (models.PaginatedResponse, error)
	GetByID(id uint) (*models.Suggestion, error)
	Update(s *models.Suggestion) error
	Store(s *models.Suggestion) error
	Delete(f *models.SuggestionFilter) error
}
