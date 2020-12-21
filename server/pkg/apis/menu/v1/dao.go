package v1

import (
	"fmt"
	rolev1 "github.com/kuops/go-example-app/server/pkg/apis/role/v1"
	"github.com/kuops/go-example-app/server/pkg/request"
	"gorm.io/gorm"
)

type dao struct {
	db *gorm.DB
}

func (s *dao)GetMenusList(model,where interface{},out interface{},pageIndex, pageSize int,totalCount *int64,whereOrder ...request.PageWhereOrder) error{
	db:=s.db.Model(model).Where(where)
	if len(whereOrder)>0 {
		for _,wo:=range whereOrder {
			if wo.Order !="" {
				db=db.Order(wo.Order)
			}
			if wo.Where !="" {
				db=db.Where(wo.Where,wo.Value...)
			}
		}
	}
	err:=db.Count(totalCount).Error
	if err!=nil{
		return err
	}
	if *totalCount==0{
		return nil
	}
	return db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}

func (s *dao)GetMenusDetail(m *Menu) (*Menu,error) {
	var menu Menu
	err := s.db.Where("id = ?",m.ID).Find(&menu).Error
	return &menu,err
}

func (s *dao)UpdateMenuInfo(m *Menu) (*Menu,error) {
	return  m,s.db.Model(m).Updates(m).Error
}

func (s *dao)CreateMenu(m *Menu) (*Menu,error) {
	err := s.db.Model(&Menu{}).Create(m).Error
	return m,err
}

func (s *dao)DeleteMenus(menuids []uint64) error {
	tx := s.db.Begin()
	for _, menuid := range menuids {
		if err := deleteMenuRecurve(tx, menuid); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("menu_id in (?)", menuids).Delete(&rolev1.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", menuids).Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func deleteMenuRecurve(db *gorm.DB, parentID uint64) error {
	var menus []Menu
	_ = db.Model(Menu{}).Where("parent_id = ?",parentID).Find(&menus)
	for _, menu := range menus {
		fmt.Println(menu.ID,menu.ParentID)
		if err := db.Where("id = ?",menu.ID).Delete(&Menu{}).Error;err != nil {
			return err
		}
		if err := db.Model(rolev1.RoleMenu{}).Where("menu_id",menu.ID).Delete(rolev1.RoleMenu{}).Error;err != nil {
			return err
		}
		if err := deleteMenuRecurve(db, menu.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *dao) GetMenuButton(adminsid uint64, menuCode string, btns *[]string) (err error) {
	sql := `select operate_type from table_menu
	      where id in (
					select menu_id from table_role_menu where 
					menu_id in (select id from table_menu where parent_id in (select id from table_menu where code=?))
					and role_id in (select role_id from table_user_role where user_id=?)
				)`
	err = s.db.Raw(sql, menuCode, adminsid).Pluck("operate_type", btns).Error
	return
}

func (s *dao)GetAllMenu() (*[]Menu,error) {
	var list []Menu
	err := s.db.Where(Menu{}).Order("parent_id asc").Order("sequence asc").Find(&list).Error
	return &list,err
}
