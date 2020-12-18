package v1

import (
	rolev1 "github.com/kuops/go-example-app/server/pkg/apis/role/v1"
	"gorm.io/gorm"
)

type dao struct {
	db *gorm.DB
}

func (s *dao)GetMenusList(limit,offset int) ([]Menu,int64,error) {
	var menus []Menu
	var total int64
	_ = s.db.Model(Menu{}).Count(&total)
	err := s.db.Model(Menu{}).Limit(limit).Offset(offset).Find(&menus).Error
	return menus,total,err
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

func (s *dao)DeleteMenus(ids []uint64) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, menuid := range ids {
		if err := deleteMenuRecurve(tx, menuid); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("menu_id in (?)", ids).Delete(&rolev1.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", ids).Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func deleteMenuRecurve(db *gorm.DB, parentID uint64) error {
	where := &Menu{}
	where.ParentID = parentID
	var menus []Menu
	dbslect := db.Where(&where)
	if err := dbslect.Find(&menus).Error; err != nil {
		return err
	}
	for _, menu := range menus {
		if err := db.Where("menu_id = ?", menu.ID).Delete(&rolev1.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := deleteMenuRecurve(db, menu.ID); err != nil {
			return err
		}
	}
	if err := dbslect.Delete(&Menu{}).Error; err != nil {
		return err
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