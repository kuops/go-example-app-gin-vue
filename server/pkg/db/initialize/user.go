package initialize

import (
	userv1 "github.com/kuops/go-example-app/server/pkg/apis/user/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"gorm.io/gorm"
	"time"
)

func newInitialUsers() *[]userv1.User {
	var users []userv1.User
	users = []userv1.User {
		{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(), UserName: "admin", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "超级管理员", HeaderImg: "https://en.gravatar.com/userimage/141081894/3772824c069f3029642247479e1664b5.jpeg", Email: "admin@example.com",Status: 1},
		{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now(), UserName: "user1", Password: "3ec063004a6f31642261936a379fde3d", NickName: "普通用户", Email: "user1@example.com",Status: 1},
	}
	return &users
}

func InitialSysUser(client *mysql.Client)  error {
	var users = newInitialUsers()
	if err := client.Database().DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []int{1, 2}).Find(&[]userv1.User{}).RowsAffected == 2 {
			log.Warn("users 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(users).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return  nil
}