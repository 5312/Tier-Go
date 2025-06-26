package model

// Role 角色模型
type Role struct {
	Base
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	DisplayName string `gorm:"size:100" json:"display_name"`
	Description string `gorm:"size:200" json:"description"`

	Users []User `gorm:"many2many:user_roles;" json:"users"`

	_ struct{} `crud:"prefix:/role,create,update,delete,page"`
}
