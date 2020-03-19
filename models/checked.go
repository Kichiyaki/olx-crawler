package models

type Checked struct {
	Model
	Value         string `gorm:"column:value" json:"value"`
	ObservationID uint   `gorm:"column:observation_id" json:"-" sql:"type:integer REFERENCES observations(id) ON DELETE CASCADE ON UPDATE RESTRICT"`
}

func (Checked) TableName() string {
	return "checked"
}
