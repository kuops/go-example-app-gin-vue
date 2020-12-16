package initialize

import (
	rolev1 "github.com/kuops/go-example-app/server/pkg/apis/role/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"gorm.io/gorm"
)

var roles = &[]rolev1.Role{
	{ID: 1, Comment: "", Name: "超级管理员", Sequence: 1, ParentID: 0},
	{ID: 2, Comment: "", Name: "普通用户", Sequence: 2, ParentID: 0},
	{ID: 3, Comment: "", Name: "普通用户子用户", Sequence: 3, ParentID: 2},
}


func InitialRole(client *mysql.Client) {
	if err := client.Database().DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []int{1, 3}).Find(&[]rolev1.Role{}).RowsAffected == 2 {
			log.Warn("table_role 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(roles).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		return nil
	}); err != nil {
		log.Panicf("table_role 表的初始数据失败,err: %v", err)
	}
}

var roleMenu = &[]rolev1.RoleMenu {
	{ID:        1, RoleID:    2, MenuID:    1},
	{ID:        2, RoleID:    2, MenuID:    2},
	{ID:        3, RoleID:    2, MenuID:    2},
	{ID:        4, RoleID:    2, MenuID:    4},
	{ID:        5, RoleID:    2, MenuID:    5},
	{ID:        6, RoleID:    2, MenuID:    6},
	{ID:        7, RoleID:    2, MenuID:    7},
	{ID:        8, RoleID:    2, MenuID:    8},
	{ID:        9, RoleID:    2, MenuID:    9},
	{ID:        10, RoleID:    2, MenuID:    10},
	{ID:        11, RoleID:    2, MenuID:    11},
	{ID:        12, RoleID:    2, MenuID:    12},
	{ID:        13, RoleID:    2, MenuID:    13},
	{ID:        14, RoleID:    2, MenuID:    14},
	{ID:        15, RoleID:    2, MenuID:    15},
	{ID:        16, RoleID:    2, MenuID:    16},
	{ID:        17, RoleID:    2, MenuID:    17},
	{ID:        18, RoleID:    2, MenuID:    18},
	{ID:        19, RoleID:    2, MenuID:    19},
	{ID:        20, RoleID:    2, MenuID:    20},
	{ID:        21, RoleID:    2, MenuID:    21},
	{ID:        22, RoleID:    2, MenuID:    22},
	{ID:        23, RoleID:    2, MenuID:    23},
}

func InitialRoleMenu(client *mysql.Client) {
	if err := client.Database().DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []int{1, 23}).Find(&[]rolev1.RoleMenu{}).RowsAffected == 2 {
			log.Warn("table_role_menu 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(roleMenu).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		return nil
	}); err != nil {
		log.Panicf("table_role_menu 表的初始数据失败,err: %v", err)
	}
}