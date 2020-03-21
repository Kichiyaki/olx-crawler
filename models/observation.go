package models

import (
	"time"
)

type Observation struct {
	Model
	Name        string    `gorm:"unique_index;column:name" json:"name"`
	URL         string    `gorm:"unique_index;column:url" json:"url"`
	Keywords    []Keyword `gorm:"foreignkey:observation_id" json:"keywords,omitempty"`
	Checked     []Checked `gorm:"foreignkey:observation_id" json:"-"`
	Started     *bool     `gorm:"column:started" json:"started,omitempty"`
	LastCheckAt time.Time `gorm:"column:last_check_at" json:"last_check_at,omitempty"`
}

func (Observation) TableName() string {
	return "observations"
}

type ObservationFilter struct {
	ID      []uint
	Name    []string
	URL     []string
	Started string
	Order   string
	Pagination
}
