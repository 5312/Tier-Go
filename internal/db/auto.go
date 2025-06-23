package db

import (
	"tier_up/app/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 系统表
		&model.User{},
		&model.Role{},
		&model.UserRole{},
	)
}
