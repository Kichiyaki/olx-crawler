package repository

import (
	"olx-crawler/errors"
	"olx-crawler/keyword"
	"olx-crawler/models"
	"olx-crawler/utils"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

type repository struct {
	db     *gorm.DB
	logrus *logrus.Entry
}

func NewKeywordRepository(db *gorm.DB) (keyword.Repository, error) {
	var err error
	for _, model := range []interface{}{
		&models.Keyword{},
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
		logrus.WithField("package", "keyword/repository"),
	}, err
}

func (repo *repository) Delete(f *models.KeywordFilter) error {
	o := []models.Keyword{}
	errs := repo.appendFilter(f).Unscoped().Delete(&o).GetErrors()
	if len(errs) > 0 {
		if f != nil {
			repo.logrus.WithField("filter", string(utils.MustMarshal(f))).Debugf("Cannot delete keywords: %v", errs)
		} else {
			repo.logrus.WithField("filter", "{}").Debugf("Cannot delete keywords: %v", errs)
		}
		return errors.Wrap(errors.ErrCannotDeleteKeywords, errs)
	}
	return nil
}

func (repo *repository) Fetch(f *models.KeywordFilter) (models.PaginatedResponse, error) {
	response := models.PaginatedResponse{}
	keywords := []*models.Keyword{}
	q := repo.appendFilter(f)
	errs := q.Find(&keywords).GetErrors()
	if len(errs) > 0 {
		if f != nil {
			repo.logrus.WithField("filter", string(utils.MustMarshal(f))).Debugf("Cannot fetch keywords: %v", errs)
		} else {
			repo.logrus.WithField("filter", "{}").Debugf("Cannot fetch keywords: %v", errs)
		}
		return response, errors.Wrap(errors.ErrCannotFetchKeywords, errs)
	}
	response.Items = keywords
	errs = q.Model(&models.Keyword{}).Limit(-1).Offset(-1).Count(&response.Total).GetErrors()
	if len(errs) > 0 {
		if f != nil {
			repo.logrus.WithField("filter", string(utils.MustMarshal(f))).Debugf("Cannot fetch keywords: %v", errs)
		} else {
			repo.logrus.WithField("filter", "{}").Debugf("Cannot fetch keywords: %v", errs)
		}
		return response, errors.Wrap(errors.ErrCannotFetchKeywords, errs)
	}
	return response, nil
}

func (repo *repository) appendFilter(f *models.KeywordFilter) *gorm.DB {
	query := repo.db
	if f != nil {
		if len(f.ID) > 0 {
			query = query.Where("id IN (?)", f.ID)
		}
		if len(f.Type) > 0 {
			query = query.Where("type IN (?)", f.Type)
		}
		if len(f.For) > 0 {
			query = query.Where("for IN (?)", f.For)
		}
		if len(f.Value) > 0 {
			query = query.Where("value IN (?)", f.Value)
		}
		if len(f.ObservationID) > 0 {
			query = query.Where("observation_id IN (?)", f.ObservationID)
		}
	}
	return query
}
