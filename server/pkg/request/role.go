package request

type DeleteRoles struct {
	IDS   []uint64 `json:"role_ids"`
}

type SetRoleMenus struct {
	RoleID  uint64 `json:"role_id"`
	MenuIDS []uint64 `json:"menu_ids"`
}


type RolePageInfo struct {
	Page     uint64 `json:"page" form:"page"`
	PageSize uint64 `json:"pageSize" form:"pageSize"`
}

type RoleList struct {
	RolePageInfo
	Key  string `json:"key"`
	Status uint64 `json:"status"`
	Sort string `json:"sort"`
}
