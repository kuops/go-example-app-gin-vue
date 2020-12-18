package v1

import (
	"fmt"
	linq "github.com/ahmetb/go-linq"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/request"
	"github.com/kuops/go-example-app/server/pkg/response"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"github.com/kuops/go-example-app/server/pkg/utils/convert"
	"github.com/kuops/go-example-app/server/pkg/utils/jwt"
	"github.com/kuops/go-example-app/server/pkg/utils/md5"
	uuid "github.com/kuops/go-example-app/server/pkg/utils/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type handler struct {
	dao *dao
	cache redis.Interface
	enforcer *casbin.Enforcer
}

func newHandler(mysqlClient *mysql.Client,redisClient redis.Interface,enforcer *casbin.Enforcer) *handler {
	return &handler{
		dao: &dao{
			db: mysqlClient.Database().DB,
		},
		cache: redisClient,
		enforcer: enforcer,
	}
}

// @Tags 登录
// @Summary 用户登录
// @Produce  application/json
// @Param data body request.Login true "用户名, 密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /api/v1/user/login [post]
func (h *handler)Login(c *gin.Context) {
		var req = request.Login{}
		_ = c.ShouldBindJSON(&req)
		if req.Username == "" || req.Password == "" {
			response.FailWithMessage("用户名密码不能为空", c)
			return
		}

		u := &User{UserName: req.Username,Password: req.Password}
		user,err := h.dao.Login(u)
		if err != nil {
			log.Errorf("用户名或密码不正确，%v",err)
			response.FailWithMessage("用户名或密码不正确", c)
			return
		}

		if user.Status != 1 {
			log.Error("用户已禁用")
			response.FailWithMessage("用户已禁用", c)
			return
		}

		key := uuid.Get()
		expires := time.Now().Add(time.Hour * 1).Unix()
		err = h.cache.Set(key,strconv.Itoa(int(user.ID)),time.Hour * 1)
		if err != nil {
			log.Errorf("redis 存储 uuid 失败",err)
			response.FailWithMessage("服务内部错误", c)
			return
		}

		cliams := &jwt.CustomClaims{
			ID: user.ID,
			UUID: key,
			Name: user.UserName,
		}
		cliams.ExpiresAt = expires
		token,err := jwt.CreateToken(cliams)
		if err != nil {
			log.Errorf("创建 Token 失败, %v",err)
			response.FailWithMessage("服务内部错误", c)
			return
		}

		resData := map[string]interface{}{
			"user": user,
			"token": token,
			"exp": expires,
		}

		response.OkWithDetailed(resData,"登录成功", c)
}

// @Tags 用户
// @Summary 用户退出
// @Produce  application/json
// @Security ApiKeyAuth
// @Success 200 {string} string "{"success":true,"data":{},"msg":"退出成功"}"
// @Router /api/v1/user/logout [post]
func (h *handler)Logout(c *gin.Context) {
	cliamsContext,_ := c.Get("claims")
	claims := cliamsContext.(*jwt.CustomClaims)
	_ = h.cache.Del(claims.UUID)
	response.OkWithMessage("退出成功",c)
}

// @Tags 用户
// @Summary 用户信息
// @Produce  application/json
// @Security ApiKeyAuth
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取用户信息成功"}"
// @Router /api/v1/user/info [get]
func (h *handler)Info(c *gin.Context) {
	cliamsContext,_ := c.Get("claims")
	claims := cliamsContext.(*jwt.CustomClaims)
	u := &User{ID: claims.ID, UserName: claims.Name}
	user,_ := h.dao.Info(u)
	var menuList []menuv1.Menu
	var err error

	if user.ID == adminID {
		menuList,err = h.dao.GetAllMenu()
		if err != nil {
			log.Errorf("获取菜单失败, %v",err)
			response.FailWithMessage("服务内部错误", c)
			return
		}
	} else {
		fmt.Println("aaa")
		menuList, err = h.dao.GetMenuByUID(user)
		fmt.Println("bbb")
		if err != nil {
			log.Errorf("获取菜单失败, %v",err)
			response.FailWithMessage("服务内部错误", c)
			return
		}
	}
	var menus  []response.MenuModel
	if len(menuList) > 0 {
		var topMenuID uint64=menuList[0].ParentID
		if topMenuID==0{
			topMenuID=menuList[0].ID
		}
		menus = setMenu(menuList, topMenuID)
	}

	resData := map[string]interface{}{
		"menus": menus,
		"user_info": user,
	}

	response.OkWithDetailed(resData,"获取用户信息成功", c)
}

func setMenu(menus []menuv1.Menu,parentID uint64) []response.MenuModel {
	var menuArr []menuv1.Menu
	var result []response.MenuModel
	linq.From(menus).Where(func(c interface{}) bool {
		return c.(menuv1.Menu).ParentID == parentID
	}).OrderBy(func(c interface{}) interface{} {
		return c.(menuv1.Menu).Sequence
	}).ToSlice(&menuArr)
	if len(menuArr) == 0 {
		return result
	}

	for _, item := range menuArr {
		menu := response.MenuModel{
			Path:      item.URL,
			Component: item.Code,
			Name:      item.Code,
			Meta:      response.MenuMeta{Title: item.Name, Icon: item.Icon},
			Children:  []response.MenuModel{}}
		if item.MenuType == 3 {
			menu.Hidden = true
		}
		//查询是否有子级
		menuChildren := setMenu(menus, item.ID)
		if len(menuChildren) > 0 {
			menu.Children = menuChildren
		}
		if item.MenuType == 2 {
			// 添加子级首页，有这一级NoCache才有效
			menuIndex := response.MenuModel{
				Path:      "index",
				Component: item.Code,
				Name:      item.Code,
				Meta:      response.MenuMeta{Title: item.Name, Icon: item.Icon},
				Children:  []response.MenuModel{}}
			menu.Children = append(menu.Children, menuIndex)
			menu.Name = menu.Name + "index"
			menu.Meta = response.MenuMeta{}
		}
		result = append(result, menu)
	}
	return result
}

// @Tags 用户
// @Summary 修改密码
// @Produce  application/json
// @Security ApiKeyAuth
// @Param data body request.ChangePassword true "用户名, 原密码, 新密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改密码成功"}"
// @Router /api/v1/user/changePassword [post]
func (h *handler)ChangePassword(c *gin.Context) {
	cliamsContext,_ := c.Get("claims")
	claims := cliamsContext.(*jwt.CustomClaims)
	u := &User{ID: claims.ID, UserName: claims.Name}
	user,_ := h.dao.Info(u)

	var req request.ChangePassword
	_ = c.ShouldBindJSON(&req)

	if req.Username == "" {
		response.FailWithMessage("用户名不能为空", c)
		return
	}
	if req.Password == "" {
		response.FailWithMessage("原始密码不能为空", c)
		return
	}
	if req.NewPassword == "" {
		response.FailWithMessage("新密码不能为空", c)
		return
	}
	if md5.Encrypt(req.Password) != user.Password {
		response.FailWithMessage("密码不正确", c)
		return
	}
	u.Password = req.NewPassword
	err := h.dao.ChangePassword(user)
	if err != nil {
		log.Errorf("修改密码失败, %v",err)
		response.FailWithMessage("修改密码失败", c)
		return
	}

	resData := map[string]interface{}{
		"user_name": user.UserName,
		"password": user.Password,
	}

	response.OkWithDetailed(resData,"修改密码成功", c)
}

// @Tags 用户
// @Summary 分页用户列表
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/v1/user/list [post]
func (h *handler)List(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	var err error

	if pageInfo.PageSize == 0 || pageInfo.Page == 0 {
		log.Error("分页参数参数非法")
		response.FailWithMessage("参数不正确", c)
		return
	}

	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	list,total,err :=  h.dao.GetUsersList(limit,offset)

	if err != nil {
		log.Errorf("获取失败 %v",err)
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}

}

// @Tags 用户
// @Summary 用户详情
// @Produce  application/json
// @Security ApiKeyAuth
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取用户详情成功"}"
// @Router /api/v1/user/detail [get]
func (h *handler)Detail(c *gin.Context) {
	cliamsContext,_ := c.Get("claims")
	claims := cliamsContext.(*jwt.CustomClaims)
	u := &User{ID: claims.ID, UserName: claims.Name}
	user,_ := h.dao.Info(u)

	resData := user
	response.OkWithDetailed(resData,"获取用户详情成功", c)

}

// @Tags 用户
// @Summary 更新用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body User true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新用户信息成功"}"
// @Router /api/v1/user/update [post]
func (h *handler) Update(c *gin.Context) {
	user := User{}
	_ = c.ShouldBindJSON(&user)
	if user.ID == 0 {
		response.FailWithMessage("用户 ID 不能为空", c)
		return
	}

	updatedUser,err := h.dao.UpdateUserInfo(&user)
	if err != nil {
		log.Errorf("更新用户信息失败 %v",err)
		response.FailWithMessage("更新用户信息失败", c)
		return
	}

	resData := map[string]interface{}{
		"user_info": updatedUser,
	}
	response.OkWithDetailed(resData, "设置成功", c)
}


// @Tags 用户
// @Summary 添加用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateUser true "用户名,密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加用户成功"}"
// @Router /api/v1/user/create [post]
func (h *handler) Create(c *gin.Context) {
	req := request.CreateUser{}
	_ = c.ShouldBindJSON(&req)
	if req.Password == "" || req.UserName == "" {
		log.Errorf("账号密码不能为空")
		response.FailWithMessage("账号密码不能为空", c)
		return
	}

	req.Password = md5.Encrypt(req.Password)
	user := User{
		UserName:    req.UserName,
		Password:    req.Password,
		NickName:    req.NickName,
		RealName:    req.RealName,
		HeaderImg:   req.HeaderImg,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Status:      req.Status,
	}
	userInfo,err := h.dao.CreateUser(&user)

	if err != nil {
		log.Errorf("添加用户失败, %v",err)
		response.FailWithMessage(fmt.Sprintf("添加用户失败, %v",err), c)
		return
	}

	resData := map[string]interface{}{
		"user_info": userInfo,
	}

	response.OkWithDetailed(resData, "添加用户成功", c)
}

// @Tags 用户
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteUsers true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除用户成功"}"
// @Router /api/v1/user/delete [post]
func (h *handler) Delete(c *gin.Context) {
	req := request.DeleteUsers{}
	var err error
	err = c.ShouldBindJSON(&req)
	if err != nil || len(req.IDS) == 0 {
		log.Errorf("删除失败,%v",err)
		response.FailWithMessage("删除失败",c)
		return
	}

	if err := h.dao.DeleteUsers(req.IDS); err != nil {
		log.Errorf("删除失败,%v",err)
		response.FailWithMessage("删除失败",c)
		return
	}

	response.OkWithMessage("删除用户成功", c)
}

// @Tags 用户
// @Summary 查看用户的角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id query int false "int valid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取用户角色成功"}"
// @Router /api/v1/user/roleList [get]
func (h *handler) RoleList (c *gin.Context) {
	var id uint64
	var err error
	var roleList []uint64
	idstr,ok := c.GetQuery("id")
	if ok {
		id,err = convert.ToUint64E(idstr)
		if err != nil {
			response.FailWithMessage("id不为数字",c)
			return
		}
	} else {
		response.FailWithMessage("id不能为空",c)
		return
	}
	userRole := UserRole{UserID: id}
	err = h.dao.GetRoleList(&UserRole{},&userRole,&roleList,"role_id")

	if err != nil {
		log.Errorf("查找 role 失败,%v",err)
		response.FailWithMessage("查找 role 失败",c)
		return
	}

	response.OkWithDetailed(&roleList,"获取用户角色成功",c)
}

// @Tags 用户
// @Summary 分配角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SetUserRole true "用户ID 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"分配角色成功"}"
// @Router /api/v1/user/setrole [post]
func (h *handler) SetRole(c *gin.Context) {
	req := request.SetUserRole{}
	_ = c.ShouldBindJSON(&req)
	err := h.dao.SetRole(req.UserID,req.RoleIDS)
	if err != nil {
		log.Errorf("分配角色失败,%v",err)
		response.FailWithMessage("分配角色失败",c)
		return
	}
	go CasbinAddRoleForUser(h.dao.db,h.enforcer,req.UserID)
	response.OkWithMessage("分配角色成功",c)
}

func CasbinAddRoleForUser(db *gorm.DB,enforcer *casbin.Enforcer,userid uint64) error {
	uid:=convert.ToString(userid)
	_,_ = enforcer.DeleteRolesForUser(uid)

	var userRoles []UserRole
	err := db.Where(&UserRole{UserID: userid}).Find(&userRoles).Error
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
	enforcer.SavePolicy()
	return nil
}
