package v1

import "time"

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserName    string       `json:"userName" gorm:"comment:用户登录名;column:username"`
	Password    string       `json:"-"  gorm:"comment:用户登录密码"`
	NickName    string       `json:"nickName" gorm:"default:系统用户;comment:用户昵称" `
	HeaderImg   string       `json:"headerImg" gorm:"default:http://qmplusimg.henrongyi.top/head.png;comment:用户头像"`
	Email       string       `json:"email" gorm:"default:demo@example.com;comment:用户邮箱"`
	RoleID      uint       	 `json:"rid" gorm:"comment:角色ID"`
}