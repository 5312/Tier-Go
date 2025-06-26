package model

// Role 角色模型
type Role struct {
	Base
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	DisplayName string `gorm:"size:100" json:"display_name"`
	Description string `gorm:"size:200" json:"description"`

	Users []User `gorm:"many2many:user_roles;" json:"-"`

	_ struct{} `crud:"prefix:/role,create,update,delete,page"`
}

type RoleReq struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}
