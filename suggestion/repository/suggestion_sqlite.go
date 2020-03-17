package repository

import (
	"olx-crawler/errors"
	"olx-crawler/models"
	"olx-crawler/suggestion"

	"github.com/jinzhu/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewSuggestionRepository(db *gorm.DB) (suggestion.Repository, error) {
	var err error
	for _, model := range []interface{}{
		&models.Suggestion{},
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

func (repo *repository) Store(s *models.Suggestion) error {
	errs := repo.db.Create(s).GetErrors()
	if len(errs) > 0 {
		return errors.Wrap(errors.ErrCannotCreateSuggestion, errs)
	}
	return nil
}

func (repo *repository) Update(input *models.Suggestion) error {
	errs := repo.
		db.
		Model(&models.Suggestion{}).
		Where("id = ?", input.ID).
		Updates(input).
		GetErrors()
	if len(errs) > 0 {
		return errors.Wrap(errors.ErrCannotUpdateSuggestion, errs)
	}
	return nil
}

func (repo *repository) Delete(f *models.SuggestionFilter) error {
	errs := repo.appendFilter(f).Delete(&[]models.Suggestion{}).GetErrors()
	if len(errs) > 0 {
		return errors.Wrap(errors.ErrCannotDeleteSuggestions, errs)
	}
	return nil
}

func (repo *repository) Fetch(f *models.SuggestionFilter) (models.PaginatedResponse, error) {
	response := models.PaginatedResponse{}
	observations := []*models.Suggestion{}
	q := repo.appendFilter(f)
	errs := q.Find(&observations).GetErrors()
	if len(errs) > 0 {
		return response, errors.Wrap(errors.ErrCannotFetchSuggestions, errs)
	}
	response.Items = observations
	errs = q.Model(&models.Suggestion{}).Limit(-1).Offset(-1).Count(&response.Total).GetErrors()
	if len(errs) > 0 {

		return response, errors.Wrap(errors.ErrCannotFetchSuggestions, errs)
	}
	return response, nil
}

func (repo *repository) GetByID(id uint) (*models.Suggestion, error) {
	o := &models.Suggestion{}
	errs := repo.db.Where("id = ?", id).First(o).GetErrors()
	if len(errs) > 0 {
		return nil, errors.Wrap(errors.ErrSuggestionNotFound, errs)
	}
	return o, nil
}

func (repo *repository) appendFilter(f *models.SuggestionFilter) *gorm.DB {
	query := repo.db
	if f != nil {
		if len(f.ID) > 0 {
			query = query.Where("id IN (?)", f.ID)
		}
		if len(f.Name) > 0 {
			query = query.Where("name IN (?)", f.Name)
		}
		if len(f.Price) > 0 {
			query = query.Where("price IN (?)", f.Price)
		}
		if len(f.ObservationID) > 0 {
			query = query.Where("observation_id IN (?)", f.ObservationID)
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
