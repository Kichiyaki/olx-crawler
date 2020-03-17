package observation

import (
	"olx-crawler/models"
)

type Usecase interface {
	Fetch(f *models.ObservationFilter) (models.PaginatedResponse, error)
	GetByID(id uint) (*models.Observation, error)
	Store(o *models.Observation) error
	Update(o *models.Observation) (*models.Observation, error)
	Delete(ids ...uint) ([]*models.Observation, error)
}
