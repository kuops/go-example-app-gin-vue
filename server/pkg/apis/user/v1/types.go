package v1

import "time"

type User struct {
	ID          uint `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint `json:"created_by" gorm:"column:created_by;default:0;not null;comment:创建者"`
	UpdatedBy   uint `json:"updated_by" gorm:"column:updated_by;default:0;not null;comment:修改者"`
	UserName    string `json:"user_name" gorm:"column:user_name;comment:用户名;"`
	Password    string `json:"-"  gorm:"column:password;comment:用户密码"`
	NickName    string `json:"nick_name" gorm:"column:nick_name;default:系统用户;comment:用户昵称" `
	RealName    string  `json:"real_name" gorm:"column:real_name;comment:真实姓名" `
	HeaderImg   string `json:"header_img" gorm:"column:header_img;default:https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif;comment:用户头像"`
	Email       string `json:"email" gorm:"column:email;default:demo@example.com;comment:用户邮箱"`
	PhoneNumber string `json:"phone_number" gorm:"column:phone_number;comment:手机号"`
	Status      int `json:"status" gorm:"column:status;comment:状态(1:正常 2:未激活 3:暂停使用)"`
}