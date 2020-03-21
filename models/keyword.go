package models

type Keyword struct {
	Model

	//excluded, one_of
	Type string `gorm:"column:type" json:"type"`
	//title, description
	For           string `gorm:"column:for" json:"for"`
	Value         string `gorm:"column:value" json:"value"`
	ObservationID uint   `gorm:"column:observation_id" json:"-" sql:"type:integer REFERENCES observations(id) ON DELETE CASCADE ON UPDATE RESTRICT"`
}

func (Keyword) TableName() string {
	return "keywords"
}

type KeywordFilter struct {
	ID            []uint
	Type          []string
	For           []string
	Value         []string
	ObservationID []uint
}
