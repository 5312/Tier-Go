package model

// User 用户模型
type User struct {
	Base

	ID       uint   `gorm:"primarykey" json:"id"`
	Username string `gorm:"size:50;not null;unique" json:"username"`
	Password string `gorm:"size:100;not null" json:"-"` // 密码不在JSON中返回
	Nickname string `gorm:"size:50" json:"nickname"`
	Email    string `gorm:"size:100;unique" json:"email"`
	Phone    string `gorm:"size:20" json:"phone"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1:正常, 0:禁用

	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

// Role 角色模型
type Role struct {
	Base
	ID uint `gorm:"primarykey" json:"id"`

	Name        string `gorm:"size:50;not null;unique" json:"name"`
	DisplayName string `gorm:"size:100" json:"display_name"`
	Description string `gorm:"size:200" json:"description"`

	Users []User `gorm:"many2many:user_roles;" json:"users"`
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}
