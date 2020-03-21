package usecase

import (
	"olx-crawler/keyword"
	"olx-crawler/models"
)

type usecase struct {
	keywordRepo keyword.Repository
}

func NewKeywordUsecase(repo keyword.Repository) keyword.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Delete(ids ...uint) ([]*models.Keyword, error) {
	f := &models.KeywordFilter{
		ID: ids,
	}
	pagination, err := ucase.keywordRepo.Fetch(f)
	if err != nil {
		return nil, err
	}
	items := pagination.Items.([]*models.Keyword)
	if err := ucase.keywordRepo.Delete(f); err != nil {
		return nil, err
	}
	return items, nil
}
