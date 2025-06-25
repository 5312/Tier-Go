package model

type Menu struct {
	Base
	ID        uint   `gorm:"primarykey" json:"id"`
	Code      uint   `json:"code"`
	Name      string `json:"name" gorm:"not null;"`
	Path      string `json:"path" gorm:"not null;comment:api路径;"`
	Component string `json:"component" gorm:"not null;comment:组件路径;" `
	ParentId  uint   `json:"parentId" gorm:"type:int(11);column:parent_id;comment:父id"`
	Icon      string `json:"icon" gorm:"comment:icon图标;"`
	Note      string `json:"note" gorm:"comment:备注;"`
	Type      int    `json:"type"`
	Status    int8   `json:"status" gorm:"comment:状态:1正常 2禁用;" `
	Sort      int    `json:"sort" gorm:"comment:显示顺序;"`
}

type Menus struct {
	Menu
	Children []Menus `json:"children"`
}
