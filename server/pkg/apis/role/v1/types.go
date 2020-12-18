package v1

import "time"


type Role struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint64 `json:"created_by" gorm:"column:created_by;default:0;not null;comment:创建者"`
	UpdatedBy   uint64 `json:"updated_by" gorm:"column:updated_by;default:0;not null;comment:修改者"`
	Comment     string `json:"comment" gorm:"column:comment;size:64;comment:备注"`
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

type Menu struct {
	ID          uint64    `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint64    `json:"created_by" gorm:"column:created_by;default:0;not null;comment:创建者"`
	UpdatedBy   uint64    `json:"updated_by" gorm:"column:updated_by;default:0;not null;comment:修改者"`
	Status      uint8     `json:"status" gorm:"column:status;type:tinyint(1);default:1;not null;comment:状态(1:启用 2:不启用)" `
	Comment     string    `json:"comment" gorm:"column:comment;size:64;comment:备注"`
	ParentID    uint64    `json:"parent_id" gorm:"column:parent_id;not null;comment:父级ID"`
	URL         string    `json:"url" gorm:"column:url;size:72;comment:菜单URL" `
	Name        string    `json:"name" gorm:"column:name;size:32;not null;comment:菜单名称" `
	Sequence    int       `json:"sequence" gorm:"column:sequence;not null;comment:排序值"`
	MenuType    uint8     `json:"menu_type" gorm:"column:menu_type;type:tinyint(1);not null;comment:菜单类型(1: 模块 2:菜单 3:操作)" `
	Code        string    `json:"code" gorm:"column:code;size:32;not null;unique_index:uk_menu_code;comment:菜单代码"`
	Icon        string    `json:"icon" gorm:"column:icon;size:32;comment:菜单图标"`
	OperateType string    `json:"operate_type" gorm:"column:operate_type;size:32;not null;comment:操作类型 none/add/del/view/update"`
}


func (RoleMenu) TableName() string {
	return "table_role_menu"
}