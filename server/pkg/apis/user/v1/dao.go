package v1

import (
	"github.com/kuops/go-example-app/server/pkg/utils/verify"
	"gorm.io/gorm"
)

type dao struct {
	db *gorm.DB
}

func (s *dao)Login(u *User) (*User,error){
	var user User
	u.Password = verify.MD5Encrypt(u.Password)
	err := s.db.Where("user_name = ? AND password = ?", u.UserName, u.Password).First(&user).Error
	return &user,err
}