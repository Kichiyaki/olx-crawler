package usecase

import (
	"olx-crawler/models"
	"olx-crawler/observation"
)

type usecase struct {
	observationRepo observation.Repository
}

func NewObservationUsecase(repo observation.Repository) observation.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Fetch(f *models.ObservationFilter) (models.PaginatedResponse, error) {
	if f == nil {
		f = &models.ObservationFilter{}
		f.Limit = 100
	}
	return ucase.observationRepo.Fetch(f)
}

func (ucase *usecase) GetByID(id uint) (*models.Observation, error) {
	return ucase.observationRepo.GetByID(id)
}

func (ucase *usecase) Store(o *models.Observation) error {
	if err := newConfig().validate(*o); err != nil {
		return err
	}

	return ucase.observationRepo.Store(o)
}

func (ucase *usecase) Update(o *models.Observation) (*models.Observation, error) {
	cfg := newConfig()
	if o.Name == "" {
		cfg.Name = false
	}
	if o.URL == "" {
		cfg.URL = false
	}
	if len(o.Exclude) == 0 {
		cfg.Exclude = false
	}
	if len(o.OneOf) == 0 {
		cfg.OneOf = false
	}
	if err := cfg.validate(*o); err != nil {
		return nil, err
	}
	if err := ucase.observationRepo.Update(o); err != nil {
		return nil, err
	}
	return ucase.GetByID(o.ID)
}

func (ucase *usecase) Delete(ids ...uint) ([]*models.Observation, error) {
	f := &models.ObservationFilter{
		ID: ids,
	}
	pagination, err := ucase.observationRepo.Fetch(f)
	if err != nil {
		return nil, err
	}
	items := pagination.Items.([]*models.Observation)
	if err := ucase.observationRepo.Delete(f); err != nil {
		return nil, err
	}
	return items, nil
}
