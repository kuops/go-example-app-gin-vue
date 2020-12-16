package v1

import "time"


type Role struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint64 `json:"created_by" gorm:"column:created_by;default:0;not null;comment:创建者"`
	UpdatedBy   uint64 `json:"updated_by" gorm:"column:updated_by;default:0;not null;comment:修改者"`
	Comment     string `json:"comment" gorm:"column:memo;size:64;comment:备注"`
	Name     string `json:"name" gorm:"column:name;size:32;not null;comment:名称"`
	Sequence int    `json:"sequence" gorm:"column:sequence;not null;comment:排序值"`
	ParentID uint64 `json:"parent_id" gorm:"column:parent_id;not null;comment:父级ID" `
}


func (Role) TableName() string {
	return "table_role"
}


type RoleMenu struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint64 `json:"created_by" gorm:"column:created_by;default:0;not null;comment:创建者"`
	UpdatedBy   uint64 `json:"updated_by" gorm:"column:updated_by;default:0;not null;comment:修改者"`
	RoleID uint64 `gorm:"column:role_id;unique_index:uk_role_menu_role_id;not null;comment:角色ID"`
	MenuID uint64 `gorm:"column:menu_id;unique_index:uk_role_menu_role_id;not null;comment:菜单ID"`
}


func (RoleMenu) TableName() string {
	return "table_role_menu"
}