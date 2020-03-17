package models

type Suggestion struct {
	Model
	Name          string       `gorm:"column:name" json:"name"`
	Price         string       `gorm:"column:price" json:"price"`
	Image         string       `gorm:"column:image" json:"image"`
	ObservationID uint         `gorm:"column:observation_id" json:"-"`
	Observation   *Observation `gorm:"foreignkey:observation_id" json:"observation,omitempty"`
}

func (Suggestion) TableName() string {
	return "suggestions"
}
