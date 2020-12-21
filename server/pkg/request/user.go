package request

type Login struct {
	Username  string `json:"username" example:"admin"`
	Password  string `json:"password" example:"123456"`
}

type ChangePassword struct {
	Username    string `json:"username" example:"admin"`
	Password    string `json:"password" example:"123456"`
	NewPassword string `json:"newPassword" example:"1234567"`
}

type UserPageInfo struct {
	Page     uint64 `json:"page" form:"page"`
	PageSize uint64 `json:"pageSize" form:"pageSize"`
}

type UserList struct {
	UserPageInfo
	Key  string `json:"key"`
	Status uint64 `json:"status"`
	Sort string `json:"sort"`
}

type PageWhereOrder struct {
	Order string
	Where string
	Value []interface{}
}

type CreateUser struct {
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	NickName    string `json:"nick_name"`
	RealName    string  `json:"real_name"`
	HeaderImg   string `json:"header_img"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Status      int `json:"status"`
}

type DeleteUsers struct {
	IDS   []uint64 `json:"user_ids"`
}

type SetUserRole struct {
	UserID uint64 `json:"user_id"`
	RoleIDS []uint64  `json:"role_ids"`
}
