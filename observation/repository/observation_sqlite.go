package repository

import (
	"olx-crawler/errors"
	"olx-crawler/models"
	"olx-crawler/observation"
	"strings"

	"github.com/jinzhu/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewObservationRepository(db *gorm.DB) (observation.Repository, error) {
	var err error
	for _, model := range []interface{}{
		&models.Observation{},
		&models.OneOf{},
		&models.Exclude{},
		&models.Checked{},
	} {
		if !db.HasTable(model) {
			errs := db.CreateTable(model).GetErrors()
			if len(errs) > 0 {
				err = errors.Wrap(errors.ErrTableCannotBeCreated, errs)
				break
			}
		}
	}
	return &repository{
		db.Set("gorm:auto_preload", true),
	}, err
}

func (repo *repository) Store(o *models.Observation) error {
	errs := repo.db.Create(o).GetErrors()
	if len(errs) > 0 {
		if strings.Contains(errs[0].Error(), "observations.url") {
			return errors.Wrap(errors.ErrObservationURLMustBeUnique, errs)
		} else if strings.Contains(errs[0].Error(), "observations.name") {
			return errors.Wrap(errors.ErrObservationNameMustBeUnique, errs)
		}
		return errors.Wrap(errors.ErrCannotCreateObservation, errs)
	}
	return nil
}

func (repo *repository) Update(input *models.Observation) error {
	errs := repo.
		db.
		Model(&models.Observation{}).
		Where("id = ?", input.ID).
		Updates(input).
		GetErrors()
	if len(errs) > 0 {
		if strings.Contains(errs[0].Error(), "observations.url") {
			return errors.Wrap(errors.ErrObservationURLMustBeUnique, errs)
		} else if strings.Contains(errs[0].Error(), "observations.name") {
			return errors.Wrap(errors.ErrObservationNameMustBeUnique, errs)
		}
		return errors.Wrap(errors.ErrCannotUpdateObservation, errs)
	}
	return nil
}

func (repo *repository) Delete(f *models.ObservationFilter) error {
	errs := repo.appendFilter(f).Delete(&[]models.Observation{}).GetErrors()
	if len(errs) > 0 {
		return errors.Wrap(errors.ErrCannotDeleteObservations, errs)
	}
	return nil
}

func (repo *repository) Fetch(f *models.ObservationFilter) (models.PaginatedResponse, error) {
	response := models.PaginatedResponse{}
	observations := []*models.Observation{}
	q := repo.appendFilter(f)
	errs := q.Find(&observations).GetErrors()
	if len(errs) > 0 {
		return response, errors.Wrap(errors.ErrCannotFetchObservations, errs)
	}
	response.Items = observations
	errs = q.Model(&models.Observation{}).Limit(-1).Offset(-1).Count(&response.Total).GetErrors()
	if len(errs) > 0 {

		return response, errors.Wrap(errors.ErrCannotFetchObservations, errs)
	}
	return response, nil
}

func (repo *repository) GetByID(id uint) (*models.Observation, error) {
	o := &models.Observation{}
	errs := repo.db.Where("id = ?", id).First(o).GetErrors()
	if len(errs) > 0 {
		return nil, errors.Wrap(errors.ErrObservationNotFound, errs)
	}
	return o, nil
}

func (repo *repository) appendFilter(f *models.ObservationFilter) *gorm.DB {
	query := repo.db
	if f != nil {
		if len(f.ID) > 0 {
			query = query.Where("id IN (?)", f.ID)
		}
		if len(f.Name) > 0 {
			query = query.Where("name IN (?)", f.Name)
		}
		if len(f.URL) > 0 {
			query = query.Where("url IN (?)", f.URL)
		}
		if f.Started != "" {
			query = query.Where("started = ?", f.Started)
		}
		if f.Order != "" {
			query = query.Order(f.Order)
		}
		if f.Limit > 0 {
			query = query.Limit(f.Limit)
		}
		if f.Offset > 0 {
			query = query.Offset(f.Offset)
		}
	}
	return query
}
