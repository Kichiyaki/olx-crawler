package models

type Checked struct {
	Model
	Value         string `gorm:"column:value" json:"value"`
	ObservationID uint   `gorm:"column:observation_id" json:"-"`
}

func (Checked) TableName() string {
	return "checked"
}
