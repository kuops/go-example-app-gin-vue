package v1

import "gorm.io/gorm"

type dao struct {
	db *gorm.DB
}

func (s *dao)GetRolesList(limit,offset int) ([]Role,int64,error) {
	var roles []Role
	var total int64
	_ = s.db.Model(Role{}).Count(&total)
	err := s.db.Model(Role{}).Limit(limit).Offset(offset).Find(&roles).Error
	return roles,total,err
}

func (s *dao)GetRolesDetail(r *Role) (*Role,error) {
	var role Role
	err := s.db.Where("id = ?",r.ID).Find(&role).Error
	return &role,err
}

func (s *dao)UpdateRoleInfo(r *Role) (*Role,error) {
	return  r,s.db.Model(r).Updates(r).Error
}

func (s *dao)CreateRole(r *Role) (*Role,error) {
	err := s.db.Model(&Role{}).Create(r).Error
	return r,err
}

func (s *dao)DeleteRoles(ids []uint64) error {
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
	if err := tx.Where("id in (?)", ids).Delete(&Role{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("role_id in (?)", ids).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *dao)GetAllRole() (*[]Role,error) {
	var list []Role
	err := s.db.Where(Role{}).Order("parent_id asc").Order("sequence asc").Find(&list).Error
	return &list,err
}


func (s *dao)GetRoleMenuIDList(model,where interface{},out interface{},fieldName string) error {
	return s.db.Model(model).Where(where).Pluck(fieldName, out).Error
}

func (s *dao)SetRoleMenus(roleid uint64, menuids []uint64) error {
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
	if err := tx.Where(&RoleMenu{RoleID: roleid}).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(menuids) > 0 {
		for _, mid := range menuids {
			rm := new(RoleMenu)
			rm.RoleID = roleid
			rm.MenuID = mid
			if err := tx.Create(rm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}