package mocks

import (
	"olx-crawler/models"
)

type Repository struct {
	FetchFunc  func(f *models.KeywordFilter) (models.PaginatedResponse, error)
	DeleteFunc func(f *models.KeywordFilter) error
}

func (m *Repository) Fetch(f *models.KeywordFilter) (models.PaginatedResponse, error) {
	return m.FetchFunc(f)
}

func (m *Repository) Delete(f *models.KeywordFilter) error {
	return m.DeleteFunc(f)
}
