package initialize

import (
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"gorm.io/gorm"
)

var menus = &[]menuv1.Menu{
	{ID: 1, ParentID: 0, Name: "TOP", Sequence: 1, MenuType: 1, Code: "TOP", Icon: "",OperateType: "none"},
	{ID: 2, ParentID: 1, Name: "系统管理", Sequence: 1, MenuType: 1, Code: "Sys", Icon: "lock", OperateType: "none",URL: "/sys"},
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
	{ID: 16, ParentID: 10, Name: "分配角色菜单", Sequence: 6, MenuType: 3, Code: "RoleSetrolemenu", Icon: "", OperateType: "setrolemenu",URL: "/role/setrole"},
	{ID: 17, ParentID: 2, Name: "后台用户管理", Sequence: 40, MenuType: 2, Code: "User", Icon: "user", OperateType: "none",URL: "/user"},
	{ID: 18, ParentID: 17, Name: "新增", Sequence: 1, MenuType: 3, Code: "UserAdd", Icon: "", OperateType: "add",URL: "/user/create"},
	{ID: 19, ParentID: 17, Name: "删除", Sequence: 2, MenuType: 3, Code: "UserDel", Icon: "", OperateType: "del",URL: "/user/delete"},
	{ID: 20, ParentID: 17, Name: "查看", Sequence: 3, MenuType: 3, Code: "UserView", Icon: "", OperateType: "view",URL: "/user/detail"},
	{ID: 21, ParentID: 17, Name: "编辑", Sequence: 4, MenuType: 3, Code: "UserUpdate", Icon: "", OperateType: "update",URL: "/user/update"},
	{ID: 22, ParentID: 17, Name: "分页api", Sequence: 5, MenuType: 3, Code: "UserList", Icon: "", OperateType: "list",URL: "/user/list"},
	{ID: 23, ParentID: 17, Name: "分配角色", Sequence: 6, MenuType: 3, Code: "UserSetrole", Icon: "", OperateType: "setuserrole",URL: "/user/setrole"},
	{ID: 24, ParentID: 1, Name: "集群管理", Sequence: 2, MenuType: 1, Code: "Cluster", Icon: "list", OperateType: "none",URL: "/cluster"},
	{ID: 25, ParentID: 24, Name: "命名空间", Sequence: 1, MenuType: 2, Code: "Namespace", Icon: "list", OperateType: "none",URL: "/cluster/namespace"},
	{ID: 26, ParentID: 25, Name: "新增", Sequence: 1, MenuType: 3, Code: "NamespaceAdd", Icon: "", OperateType: "add",URL: "/cluster/namespace/create"},
	{ID: 27, ParentID: 25, Name: "删除", Sequence: 2, MenuType: 3, Code: "NamespaceDel", Icon: "", OperateType: "del",URL: "/cluster/namespace/delete"},
	{ID: 28, ParentID: 25, Name: "查看", Sequence: 3, MenuType: 3, Code: "NamespaceView", Icon: "", OperateType: "view",URL: "/cluster/namespace/detail"},
	{ID: 29, ParentID: 25, Name: "编辑", Sequence: 4, MenuType: 3, Code: "NamespaceUpdate", Icon: "", OperateType: "update",URL: "/cluster/namespaceupdate"},
	{ID: 30, ParentID: 25, Name: "分页api", Sequence: 5, MenuType: 3, Code: "NamespaceList", Icon: "", OperateType: "list",URL: "/cluster/namespace/list"},
	{ID: 31, ParentID: 24, Name: "无状态部署", Sequence: 1, MenuType: 2, Code: "Deployment", Icon: "list", OperateType: "none",URL: "/cluster/deployment"},
	{ID: 32, ParentID: 31, Name: "新增", Sequence: 1, MenuType: 3, Code: "DeploymentAdd", Icon: "", OperateType: "add",URL: "/cluster/deployment/create"},
	{ID: 33, ParentID: 31, Name: "删除", Sequence: 2, MenuType: 3, Code: "DeploymentDel", Icon: "", OperateType: "del",URL: "/cluster/deployment/delete"},
	{ID: 34, ParentID: 31, Name: "查看", Sequence: 3, MenuType: 3, Code: "DeploymentView", Icon: "", OperateType: "view",URL: "/cluster/deployment/detail"},
	{ID: 35, ParentID: 31, Name: "编辑", Sequence: 4, MenuType: 3, Code: "DeploymentUpdate", Icon: "", OperateType: "update",URL: "/cluster/deployment/update"},
	{ID: 36, ParentID: 31, Name: "编辑", Sequence: 5, MenuType: 3, Code: "DeploymentList", Icon: "", OperateType: "update",URL: "/cluster/deployment/list"},
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
