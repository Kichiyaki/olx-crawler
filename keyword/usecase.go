package keyword

import (
	"olx-crawler/models"
)

type Usecase interface {
	Delete(ids ...uint) ([]*models.Keyword, error)
}
