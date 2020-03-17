package models

type Exclude struct {
	Model

	//title, details or description
	For           string `gorm:"column:for" json:"for"`
	Value         string `gorm:"column:value" json:"value"`
	ObservationID uint   `gorm:"column:observation_id" json:"-"`
}

func (Exclude) TableName() string {
	return "exclude"
}
