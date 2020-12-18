package request

type DeleteRoles struct {
	IDS   []uint64 `json:"user_ids"`
}

type SetRoleMenus struct {
	RoleID  uint64 `json:"role_id"`
	MenuIDS []uint64 `json:"menu_ids"`
}
