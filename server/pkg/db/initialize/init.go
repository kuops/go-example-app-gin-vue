package initialize

import (
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	rolev1 "github.com/kuops/go-example-app/server/pkg/apis/role/v1"
	userv1 "github.com/kuops/go-example-app/server/pkg/apis/user/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
)

func Migration(client *mysql.Client) {
	err := client.Database().DB.AutoMigrate(
		userv1.User{},
		userv1.UserRole{},
		menuv1.Menu{},
		rolev1.Role{},
		rolev1.RoleMenu{},
		)
	if err != nil {
		log.Panicf("初始化表结构失败")
	}
}

func InitialDatas(client *mysql.Client) {
	InitialUser(client)
	InitialRole(client)
	InitialUserRole(client)
	InitialMenu(client)
	InitialRoleMenu(client)
}
