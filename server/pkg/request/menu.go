package request

type DeleteMenus struct {
	IDS   []uint64 `json:"menu_ids"`
}

type MenuPageInfo struct {
	Page     uint64 `json:"page" form:"page"`
	PageSize uint64 `json:"pageSize" form:"pageSize"`
}

type MenuList struct {
	RolePageInfo
	Key  string `json:"key"`
	Sort string `json:"sort"`
	MenuType uint64 `json:"type"`
	ParentID uint64 `json:"parent_id"`
}
