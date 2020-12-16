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

type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
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