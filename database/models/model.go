package models

import (
	"time"
)

// Model struct, every struct models will have this type
type Model struct {
	ID        uint64     `gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
}
