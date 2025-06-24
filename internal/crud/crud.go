package crud

import (
	"gorm.io/gorm"
)

type ICrud[T any] interface {
	Create(*T) error
	// FindByID(id uint) (*T, error)
	// Update(*T) error
	// Delete(id uint) error
	// Page() error
}

type Crud[T any] struct {
	DB *gorm.DB
}

func (c *Crud[T]) Create(entity *T) error {
	if err := c.DB.Create(entity).Error; err != nil {
		return err
	}
	return nil
}
