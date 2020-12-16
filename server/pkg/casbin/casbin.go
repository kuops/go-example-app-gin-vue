package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	rolev1 "github.com/kuops/go-example-app/server/pkg/apis/role/v1"
	userv1 "github.com/kuops/go-example-app/server/pkg/apis/user/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/utils/convert"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer
var db *gorm.DB

// 角色-URL导入
func InitCsbinEnforcer(client *mysql.Client) {
	db = client.Database().DB
	var enforcer *casbin.Enforcer

	casbinModel := `# 用于request的定义，sub, obj, act 表示经典三元组: 访问实体 (Subject)，访问资源 (Object) 和访问方法 (Action)
	[request_definition]
	r = sub, obj, act

	# 对policy的定义，sub, obj, act 表示经典三元组: 访问实体 (Subject)，访问资源 (Object) 和访问方法 (Action)
	[policy_definition]
	p = sub, obj, act
	
    # g 是一个 RBAC系统，_, _表示角色继承关系的前项和后项，一般来讲，如果您需要进行角色和用户的绑定，直接使用g 即可
	[role_definition]
	g = _, _
	
    # 一个只有一条规则生效，其余都被拒绝的情况，
    # 该Effect原语表示如果存在任意一个决策结果为allow的匹配规则，则最终决策结果为allow，
    # 即allow-override。 其中p.eft 表示策略规则的决策结果，可以为allow 或者deny
    # 当不指定规则的决策结果时,取默认值allow 。 通常情况下，policy的p.eft默认为allow， 因此前面例子中都使用了这个默认值
	[policy_effect]
	e = some(where (p.eft == allow))

    # 表示请求的三元组：主题、对象、行为都应该匹配策略规则中的表达式。
    # keyMatch2 只要 key match 就可以了
    # 
	[matchers]
    m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)`

	m,err := model.NewModelFromString(casbinModel)
	if err != nil {
		println(err.Error())
	}
	conn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		client.Config.Username,
		client.Config.Password,
		client.Config.Host,
		client.Config.Port,
		client.Config.Database)
	a,_ := gormadapter.NewAdapter("mysql" , conn,true)
	enforcer, err = casbin.NewEnforcer(m,a)
	if err != nil {
		println(err)
	}
	_= enforcer.LoadPolicy()

	var roles []rolev1.Role
	err = db.Where(&rolev1.Role{}).Find(&roles).Error
	if err != nil {
		log.Panic("err: %v",err)
	}

	for _, role := range roles {
		setRolePermission(enforcer, role.ID)
	}

	var users []userv1.User
	err = db.Where(&userv1.User{}).Find(&users).Error
	if err != nil {
		log.Panic("err: %v",err)
	}

	for _, user := range users {
		if err = CsbinAddRoleForUser(enforcer, user.ID);err != nil {
			log.Panic("err: %v",err)
		}
	}

	_ = enforcer.SavePolicy()
	Enforcer = enforcer
}

func setRolePermission(enforcer *casbin.Enforcer, roleid uint64) {
	var roleMenus []rolev1.RoleMenu
	err := db.Where(&rolev1.RoleMenu{RoleID: roleid}).Find(&roleMenus).Error
	if err != nil {
		return
	}
	_,_ = enforcer.AddPermissionForUser(convert.ToString(roleid), "/api/v1/user/info","GET|POST|UPDATE|DELETE")
	_,_ = 	enforcer.AddPermissionForUser(convert.ToString(roleid), "/api/v1/user/changePassword","GET|POST|UPDATE|DELETE")
	for _, roleMenu := range roleMenus {
		menu := menuv1.Menu{}
		where := menuv1.Menu{}
		where.ID = roleMenu.MenuID
		err = db.Where(&where).First(&menu).Error
		if err != nil {
			return
		}
		if  menu.MenuType == 3 {
			_,err := enforcer.AddPermissionForUser(convert.ToString(roleid), "/api/v1" + menu.URL,"GET|POST|UPDATE|DELETE")
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

// 用户角色处理
func CsbinAddRoleForUser(enforcer *casbin.Enforcer,userid uint64) error {
	uid:=convert.ToString(userid)
	_,_ = enforcer.DeleteRolesForUser(uid)

	var userRoles []userv1.UserRole
	err := db.Where(&userv1.UserRole{UserID: userid}).Find(&userRoles).Error
	if err != nil {
		return err
	}

	roles := []string{}
	for _, userRole := range userRoles {
		roles = append(roles, convert.ToString(userRole.RoleID))
	}

	_,err = enforcer.AddRolesForUser(uid, roles)
	if err != nil {
		panic(err)
	}
	return nil
}