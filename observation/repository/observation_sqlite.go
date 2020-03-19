package repository

import (
	"olx-crawler/errors"
	"olx-crawler/models"
	"olx-crawler/observation"
	"olx-crawler/utils"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

type repository struct {
	db     *gorm.DB
	logrus *logrus.Entry
}

func NewObservationRepository(db *gorm.DB) (observation.Repository, error) {
	var err error
	for _, model := range []interface{}{
		&models.Observation{},
		&models.OneOf{},
		&models.Excluded{},
		&models.Checked{},
	} {
		if !db.HasTable(model) {
			errs := db.CreateTable(model).GetErrors()
			if len(errs) > 0 {
				logrus.Debugf("Cannot create table: %v", errs)
				err = errors.Wrap(errors.ErrTableCannotBeCreated, errs)
				break
			}
		}
	}
	return &repository{
		db.Set("gorm:auto_preload", true),
		logrus.WithField("package", "observation/repository"),
	}, err
}

func (repo *repository) Store(o *models.Observation) error {
	errs := repo.db.Create(o).GetErrors()
	if len(errs) > 0 {
		repo.logrus.WithField("observation", string(utils.MustMarshal(o))).Debugf("Cannot store observation: %v", errs)
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
		repo.logrus.WithField("observation", string(utils.MustMarshal(input))).Debugf("Cannot update observation: %v", errs)
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
	o := []models.Observation{}
	errs := repo.appendFilter(f).Unscoped().Delete(&o).GetErrors()
	if len(errs) > 0 {
		if f != nil {
			repo.logrus.WithField("filter", string(utils.MustMarshal(f))).Debugf("Cannot delete observations: %v", errs)
		} else {
			repo.logrus.WithField("filter", "{}").Debugf("Cannot delete observations: %v", errs)
		}
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
		if f != nil {
			repo.logrus.WithField("filter", string(utils.MustMarshal(f))).Debugf("Cannot fetch observations: %v", errs)
		} else {
			repo.logrus.WithField("filter", "{}").Debugf("Cannot fetch observations: %v", errs)
		}
		return response, errors.Wrap(errors.ErrCannotFetchObservations, errs)
	}
	response.Items = observations
	errs = q.Model(&models.Observation{}).Limit(-1).Offset(-1).Count(&response.Total).GetErrors()
	if len(errs) > 0 {
		if f != nil {
			repo.logrus.WithField("filter", string(utils.MustMarshal(f))).Debugf("Cannot fetch observations: %v", errs)
		} else {
			repo.logrus.WithField("filter", "{}").Debugf("Cannot fetch observations: %v", errs)
		}
		return response, errors.Wrap(errors.ErrCannotFetchObservations, errs)
	}
	return response, nil
}

func (repo *repository) GetByID(id uint) (*models.Observation, error) {
	o := &models.Observation{}
	errs := repo.db.Where("id = ?", id).First(o).GetErrors()
	if len(errs) > 0 {
		repo.logrus.WithField("id", id).Debugf("Cannot get observation: %v", errs)
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
		if f.Started == "true" {
			query = query.Where("started = true")
		} else if f.Started == "false" {
			query = query.Where("started = false")
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
