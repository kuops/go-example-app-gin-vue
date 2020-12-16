package initialize

import (
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"gorm.io/gorm"
)

var menus = &[]menuv1.Menu{
	{ID: 1, ParentID: 0, Name: "TOP", Sequence: 1, MenuType: 1, Code: "TOP", Icon: "",OperateType: "none"},
	{ID: 2, ParentID: 1, Name: "系统管理", Sequence: 1, MenuType: 1, Code: "Sys", Icon: "lock", OperateType: "none"},
	{ID: 3, ParentID: 2, Name: "图标管理", Sequence: 10, MenuType: 2, Code: "Icon", Icon: "icon", OperateType: "none",URL: "/icon"},
	{ID: 4, ParentID: 2, Name: "菜单管理", Sequence: 20, MenuType: 2, Code: "Menu", Icon: "documentation", OperateType: "none",URL: "/menu"},
	{ID: 5, ParentID: 4, Name: "新增", Sequence: 1, MenuType: 3, Code: "MenuAdd", Icon: "", OperateType: "add",URL: "/menu/create"},
	{ID: 6, ParentID: 4, Name: "删除", Sequence: 2, MenuType: 3, Code: "MenuDel", Icon: "", OperateType: "del",URL: "/menu/delete"},
	{ID: 7, ParentID: 4, Name: "查看", Sequence: 3, MenuType: 3, Code: "MenuView", Icon: "", OperateType: "view",URL: "/menu/detail"},
	{ID: 8, ParentID: 4, Name: "编辑", Sequence: 4, MenuType: 3, Code: "MenuUpdate", Icon: "", OperateType: "update",URL: "/menu/update"},
	{ID: 9, ParentID: 4, Name: "分页api", Sequence: 5, MenuType: 3, Code: "MenuList", Icon: "", OperateType: "list",URL: "/menu/list"},
	{ID: 10, ParentID: 2, Name: "角色管理", Sequence: 30, MenuType: 2, Code: "Role", Icon: "tree", OperateType: "none",URL: "/role"},
	{ID: 11, ParentID: 10, Name: "新增", Sequence: 1, MenuType: 3, Code: "RoleAdd", Icon: "", OperateType: "add",URL: "/role/create"},
	{ID: 12, ParentID: 10, Name: "删除", Sequence: 2, MenuType: 3, Code: "RoleDel", Icon: "", OperateType: "del",URL: "/role/delete"},
	{ID: 13, ParentID: 10, Name: "查看", Sequence: 3, MenuType: 3, Code: "RoleView", Icon: "", OperateType: "view",URL: "/role/detail"},
	{ID: 14, ParentID: 10, Name: "编辑", Sequence: 4, MenuType: 3, Code: "RoleUpdate", Icon: "", OperateType: "update",URL: "/role/update"},
	{ID: 15, ParentID: 10, Name: "分页api", Sequence: 5, MenuType: 3, Code: "RoleList", Icon: "", OperateType: "list",URL: "/role/list"},
	{ID: 16, ParentID: 10, Name: "分配角色菜单", Sequence: 6, MenuType: 3, Code: "RoleSetrolemenu", Icon: "", OperateType: "set",URL: "/role/menu"},
	{ID: 17, ParentID: 2, Name: "后台用户管理", Sequence: 40, MenuType: 2, Code: "Users", Icon: "user", OperateType: "none",URL: "/user"},
	{ID: 18, ParentID: 17, Name: "新增", Sequence: 1, MenuType: 3, Code: "UserAdd", Icon: "", OperateType: "add",URL: "/user/create"},
	{ID: 19, ParentID: 17, Name: "删除", Sequence: 2, MenuType: 3, Code: "UserDel", Icon: "", OperateType: "del",URL: "/user/delete"},
	{ID: 20, ParentID: 17, Name: "查看", Sequence: 3, MenuType: 3, Code: "UserView", Icon: "", OperateType: "view",URL: "/user/detail"},
	{ID: 21, ParentID: 17, Name: "编辑", Sequence: 4, MenuType: 3, Code: "UserUpdate", Icon: "", OperateType: "update",URL: "/user/update"},
	{ID: 22, ParentID: 17, Name: "分页api", Sequence: 5, MenuType: 3, Code: "UserList", Icon: "", OperateType: "list",URL: "/user/list"},
	{ID: 23, ParentID: 17, Name: "分配角色", Sequence: 6, MenuType: 3, Code: "UserSetrole", Icon: "", OperateType: "set",URL: "/user/role"},
}


func InitialMenu(client *mysql.Client) {
	if err := client.Database().DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []int{1, 23}).Find(&[]menuv1.Menu{}).RowsAffected == 2 {
			log.Warn("table_menu 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(menus).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		return nil
	}); err != nil {
		log.Panicf("table_menu 表的初始数据失败,err: %v", err)
	}
}