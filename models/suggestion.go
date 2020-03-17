package models

type Suggestion struct {
	Model
	Name          string       `gorm:"column:name" json:"name"`
	Price         string       `gorm:"column:price" json:"price"`
	Image         string       `gorm:"column:image" json:"image"`
	URL           string       `gorm:"column:url" json:"url"`
	ObservationID uint         `gorm:"column:observation_id" json:"-"`
	Observation   *Observation `gorm:"foreignkey:observation_id" json:"observation,omitempty"`
}

func (Suggestion) TableName() string {
	return "suggestions"
}

type SuggestionFilter struct {
	ID            []uint
	Name          []string
	Price         []string
	ObservationID []uint
	Order         string
	Pagination
}
