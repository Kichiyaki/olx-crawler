package keyword

import (
	"olx-crawler/models"
)

type Repository interface {
	Fetch(f *models.KeywordFilter) (models.PaginatedResponse, error)
	Delete(f *models.KeywordFilter) error
}
