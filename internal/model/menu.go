package model

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	Id        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Code string ``
}
