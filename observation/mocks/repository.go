package mocks

import (
	"olx-crawler/models"
)

type Repository struct {
	FetchFunc   func(f *models.ObservationFilter) (models.PaginatedResponse, error)
	GetByIDFunc func(id uint) (*models.Observation, error)
	StoreFunc   func(o *models.Observation) error
	UpdateFunc  func(o *models.Observation) error
	DeleteFunc  func(f *models.ObservationFilter) error
}

func (m *Repository) Fetch(f *models.ObservationFilter) (models.PaginatedResponse, error) {
	return m.FetchFunc(f)
}

func (m *Repository) GetByID(id uint) (*models.Observation, error) {
	return m.GetByIDFunc(id)
}

func (m *Repository) Store(o *models.Observation) error {
	return m.StoreFunc(o)
}

func (m *Repository) Update(o *models.Observation) error {
	return m.UpdateFunc(o)
}

func (m *Repository) Delete(f *models.ObservationFilter) error {
	return m.DeleteFunc(f)
}
