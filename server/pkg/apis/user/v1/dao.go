package v1

import (
	"errors"
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	"github.com/kuops/go-example-app/server/pkg/request"
	"github.com/kuops/go-example-app/server/pkg/utils/md5"
	"gorm.io/gorm"
)

type dao struct {
	db *gorm.DB
}

func (s *dao)Login(u *User) (*User,error){
	var user User
	u.Password = md5.Encrypt(u.Password)
	err := s.db.Where("user_name = ? AND password = ?", u.UserName, u.Password).First(&user).Error
	return &user,err
}

func (s *dao)Info(u *User) (*User,error){
	var user User
	err := s.db.Where("id = ?", u.ID).First(&user).Error
	return &user,err
}

func (s *dao)GetAllMenu() ([]menuv1.Menu,error){
	var menu []menuv1.Menu
	err := s.db.Where(&menuv1.Menu{}).Find(&menu).Error
	return menu,err
}

func (s *dao)GetMenuByUID(u *User) ([]menuv1.Menu,error){
	var menu []menuv1.Menu
	sql := `select * FROM table_menu 
	WHERE id in (select menu_id from table_role_menu WHERE role_id in (SELECT role_id FROM table_user_role WHERE user_id = ?))`
	err := s.db.Debug().Raw(sql, u.ID).Find(&menu).Error
	return menu,err
}

func (s *dao)ChangePassword(u *User) error {
	password := md5.Encrypt(u.Password)
	err := s.db.Debug().Model(User{}).Where("id = ? AND user_name = ?", u.ID, u.UserName).Update("password",password).Error
	return err
}

func (s *dao)GetUsersList(model,where interface{},out interface{},pageIndex, pageSize int,totalCount *int64,whereOrder ...request.PageWhereOrder) error{
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

func (s *dao)UpdateUserInfo(u *User) (*User,error) {
	err := s.db.Model(u).Updates(u).Error
	return u,err
}

func (s *dao)CreateUser(u *User) (*User,error) {
	var user User
	if !errors.Is(s.db.Where("user_name = ?", u.UserName).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return &User{},errors.New("用户名已注册")
	}
	err := s.db.Model(&User{}).Create(u).Error
	return u,err
}


func (s *dao)DeleteUsers(ids []uint64) error {
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
	if err := tx.Where("id in (?)", ids).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("user_id in (?)", ids).Delete(&UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *dao)GetRoleList(model,where interface{},out interface{},fieldName string) error {
	return s.db.Model(model).Where(where).Pluck(fieldName, out).Error
}

func (s *dao) SetRole(adminsid uint64, roleids []uint64) error {
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
	if err := tx.Where(&UserRole{UserID: adminsid}).Delete(&UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(roleids) > 0 {
		for _, rid := range roleids {
			rm := new(UserRole)
			rm.RoleID = rid
			rm.UserID = adminsid
			if err := tx.Create(rm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}
