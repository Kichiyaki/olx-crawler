package models

type Suggestion struct {
	Model
	OlxID         string       `gorm:"column:olx_id" json:"olx_id"`
	Title         string       `gorm:"column:title" json:"title"`
	Price         string       `gorm:"column:price" json:"price"`
	Image         string       `gorm:"column:image" json:"image"`
	URL           string       `gorm:"column:url" json:"url"`
	ObservationID uint         `gorm:"column:observation_id" json:"-" sql:"type:integer REFERENCES observations(id) ON DELETE CASCADE ON UPDATE RESTRICT"`
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
