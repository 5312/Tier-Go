package model

type Menu struct {
	Base
	Id   uint   `gorm:"primarykey" json:"id"`
	Code string ``
}
