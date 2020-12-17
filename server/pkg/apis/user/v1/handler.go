package v1

import (
	"fmt"
	linq "github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/request"
	"github.com/kuops/go-example-app/server/pkg/response"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"github.com/kuops/go-example-app/server/pkg/utils/jwt"
	"github.com/kuops/go-example-app/server/pkg/utils/md5"
	uuid "github.com/kuops/go-example-app/server/pkg/utils/uuid"
	"strconv"
	"time"
)

type handler struct {
	dao *dao
	cache redis.Interface
}

func newHandler(mysqlClient *mysql.Client,redisClient redis.Interface) *handler {
	return &handler{
		dao: &dao{
			db: mysqlClient.Database().DB,
		},
		cache: redisClient,
	}
}

// @Tags 用户
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
	var menus  []MenuModel
	if len(menuList) > 0 {
		var topMenuID uint64=menuList[0].ParentID
		if topMenuID==0{
			topMenuID=menuList[0].ID
		}
		menus = setMenu(menuList, topMenuID)
	}

	if len(menus) == 0 && user.ID == adminID {
		menus = getSuperAdminMenu()
	}

	resData := map[string]interface{}{
		"menus": menus,
		"user_info": user,
	}

	response.OkWithDetailed(resData,"获取用户信息成功", c)
}

func getSuperAdminMenu() (out []MenuModel) {
	menuTop := MenuModel{
		Path:      "/sys",
		Component: "Sys",
		Name:      "Sys",
		Meta:      MenuMeta{Title: "系统管理", NoCache: false},
		Children:  []MenuModel{}}
	menuModel := MenuModel{
		Path:      "/icon",
		Component: "Icon",
		Name:      "Icon",
		Meta:      MenuMeta{Title: "图标管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/menu",
		Component: "Menu",
		Name:      "Menu",
		Meta:      MenuMeta{Title: "菜单管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/role",
		Component: "Role",
		Name:      "Role",
		Meta:      MenuMeta{Title: "角色管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/user",
		Component: "Admins",
		Name:      "Admins",
		Meta:      MenuMeta{Title: "用户管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	out = append(out, menuTop)
	return
}

func setMenu(menus []menuv1.Menu,parentID uint64) []MenuModel {
	var menuArr []menuv1.Menu
	var result []MenuModel
	linq.From(menus).Where(func(c interface{}) bool {
		return c.(menuv1.Menu).ParentID == parentID
	}).OrderBy(func(c interface{}) interface{} {
		return c.(menuv1.Menu).Sequence
	}).ToSlice(&menuArr)
	if len(menuArr) == 0 {
		return result
	}
	noCache := false
	for _, item := range menuArr {
		menu := MenuModel{
			Path:      item.URL,
			Component: item.Code,
			Name:      item.Code,
			Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
			Children:  []MenuModel{}}
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
			menuIndex := MenuModel{
				Path:      "index",
				Component: item.Code,
				Name:      item.Code,
				Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
				Children:  []MenuModel{}}
			menu.Children = append(menu.Children, menuIndex)
			menu.Name = menu.Name + "index"
			menu.Meta = MenuMeta{}
		}
		result = append(result, menu)
	}
	return result
}

type MenuMeta struct {
	Title   string `json:"title"`   // 标题
	Icon    string `json:"icon"`    // 图标
	NoCache bool   `json:"noCache"` // 是不是缓存
}

type MenuModel struct {
	Path      string      `json:"path"`      // 路由
	Component string      `json:"component"` // 对应vue中的map name
	Name      string      `json:"name"`      // 菜单名称
	Hidden    bool        `json:"hidden"`    // 是否隐藏
	Meta      MenuMeta    `json:"meta"`      // 菜单信息
	Children  []MenuModel `json:"children"`  // 子级菜单
}

type UserData struct {
	Menus        []MenuModel `json:"menus"`        // 菜单
	Introduction string      `json:"introduction"` // 介绍
	Avatar       string      `json:"avatar"`       // 图标
	Name         string      `json:"name"`         // 姓名
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
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改密码成功"}"
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
	list,total,err :=  h.dao.GetUserInfoList(limit,offset)

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

	resData := map[string]interface{}{
		"user_info": user,
	}

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