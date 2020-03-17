package models

type OneOf struct {
	Model

	//title, details or description
	For           string `gorm:"column:for" json:"for"`
	Value         string `gorm:"column:value" json:"value"`
	ObservationID uint   `gorm:"column:observation_id" json:"-"`
}

func (OneOf) TableName() string {
	return "one_of"
}
