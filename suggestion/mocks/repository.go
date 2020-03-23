package mocks

import (
	"olx-crawler/models"
)

type Repository struct {
	FetchFunc   func(f *models.SuggestionFilter) (models.PaginatedResponse, error)
	GetByIDFunc func(id uint) (*models.Suggestion, error)
	StoreFunc   func(o *models.Suggestion) error
	UpdateFunc  func(o *models.Suggestion) error
	DeleteFunc  func(f *models.SuggestionFilter) error
}

func (m *Repository) Fetch(f *models.SuggestionFilter) (models.PaginatedResponse, error) {
	return m.FetchFunc(f)
}

func (m *Repository) GetByID(id uint) (*models.Suggestion, error) {
	return m.GetByIDFunc(id)
}

func (m *Repository) Store(o *models.Suggestion) error {
	return m.StoreFunc(o)
}

func (m *Repository) Update(o *models.Suggestion) error {
	return m.UpdateFunc(o)
}

func (m *Repository) Delete(f *models.SuggestionFilter) error {
	return m.DeleteFunc(f)
}
