package observation

import (
	"olx-crawler/models"
)

type Repository interface {
	Fetch(f *models.ObservationFilter) (models.PaginatedResponse, error)
	GetByID(id uint) (*models.Observation, error)
	Update(o *models.Observation) error
	Store(o *models.Observation) error
	Delete(f *models.ObservationFilter) error
}
