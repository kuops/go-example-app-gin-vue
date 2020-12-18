package v1

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/request"
	"github.com/kuops/go-example-app/server/pkg/response"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"github.com/kuops/go-example-app/server/pkg/utils/convert"
	"github.com/kuops/go-example-app/server/pkg/utils/jwt"
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


// @Tags 菜单
// @Summary 分页菜单列表
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/v1/menu/list [post]
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
	list,total,err :=  h.dao.GetMenusList(limit,offset)

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


// @Tags 菜单
// @Summary 菜单详情
// @Produce  application/json
// @Security ApiKeyAuth
// @Param id query int false "int valid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取用户详情成功"}"
// @Router /api/v1/menu/detail [get]
func (h *handler)Detail(c *gin.Context) {
	var id uint64
	var err error
	var menu Menu
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
	menu.ID = id
	m,err := h.dao.GetMenusDetail(&menu)
	if err != nil {
		log.Errorf("获取失败 %v",err)
		response.FailWithMessage("获取失败", c)
	}

	resData := map[string]interface{}{
		"menu_detail": m,
	}

	response.OkWithDetailed(resData, "获取成功", c)
}

// @Tags 菜单
// @Summary 更新菜单信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body Menu true "ID, 父级ID, URL"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新菜单信息成功"}"
// @Router /api/v1/menu/update [post]
func (h *handler) Update(c *gin.Context) {
	menu := Menu{}
	_ = c.ShouldBindJSON(&menu)
	if menu.ID == 0 {
		response.FailWithMessage("菜单 ID 不能为空", c)
		return
	}

	updatedMenu,err := h.dao.UpdateMenuInfo(&menu)
	if err != nil {
		log.Errorf("更新菜单信息失败 %v",err)
		response.FailWithMessage("更新菜单信息失败", c)
		return
	}

	resData := map[string]interface{}{
		"menu_info": updatedMenu,
	}

	response.OkWithDetailed(resData, "更新菜单信息成功", c)
}

// @Tags 菜单
// @Summary 删除菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteMenus true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除用户成功"}"
// @Router /api/v1/user/delete [post]
func (h *handler) Delete(c *gin.Context) {
	req := request.DeleteMenus{}
	var err error
	err = c.ShouldBindJSON(&req)
	if err != nil || len(req.IDS) == 0 {
		log.Errorf("删除失败,%v",err)
		response.FailWithMessage("删除失败",c)
		return
	}

	if err := h.dao.DeleteMenus(req.IDS); err != nil {
		log.Errorf("删除失败,%v",err)
		response.FailWithMessage("删除失败",c)
		return
	}

	response.OkWithMessage("删除用户成功", c)
}


// @Tags 菜单
// @Summary 添加菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body Menu true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加菜单成功"}"
// @Router /api/v1/menu/create [post]
func (h *handler) Create(c *gin.Context) {
	req := Menu{}
	err := c.ShouldBindJSON(&req)

	if err != nil {
		response.FailWithMessage("缺少必填项", c)
		return
	}

	m,err := h.dao.CreateMenu(&req)

	if err != nil {
		log.Errorf("添加菜单失败, %v",err)
		response.FailWithMessage("添加菜单失败", c)
		return
	}

	resData := map[string]interface{}{
		"menu_info": m,
	}

	go h.InitMenu(m)
	response.OkWithDetailed(resData, "添加菜单成功", c)
}

// 新增菜单后自动添加菜单下的常规操作
func  (h *handler)  InitMenu(model *Menu) {
	if model.MenuType != 2 {
		return
	}

	add := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/create", Name: "新增", Sequence: 1, MenuType: 3, Code: model.Code + "Add", OperateType: "add"}
	_,_ = h.dao.CreateMenu(&add)
	del := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/delete", Name: "删除", Sequence: 2, MenuType: 3, Code: model.Code + "Del", OperateType: "del"}
	_,_ = h.dao.CreateMenu(&del)
	view := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/detail", Name: "查看", Sequence: 3, MenuType: 3, Code: model.Code + "View", OperateType: "view"}
	_,_ = h.dao.CreateMenu(&view)
	update := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/update", Name: "编辑", Sequence: 4, MenuType: 3, Code: model.Code + "Update", OperateType: "update"}
	_,_ = h.dao.CreateMenu(&update)
	list := Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/list", Name: "分页api", Sequence: 5, MenuType: 3, Code: model.Code + "List", OperateType: "list"}
	_,_ = h.dao.CreateMenu(&list)
}


// @Tags 菜单
// @Summary 获取有权限操作的菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param menucode query string false "int valid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加菜单成功"}"
// @Router /api/v1/menu/menubuttonlist [get]
func (h *handler) MenuButtonList(c *gin.Context) {
	cliamsContext,_ := c.Get("claims")
	claims := cliamsContext.(*jwt.CustomClaims)
	menuCode,_ := c.GetQuery("menucode")

	if claims.ID == 0 || menuCode == "" {
		response.FailWithMessage("获取操作列表失败", c)
		return
	}

	var btnList []string
	if claims.ID == 1 {
		//管理员
		btnList = append(btnList, "add")
		btnList = append(btnList, "del")
		btnList = append(btnList, "view")
		btnList = append(btnList, "update")
		btnList = append(btnList, "setrolemenu")
		btnList = append(btnList, "setuserrole")
	} else {
		err := h.dao.GetMenuButton(claims.ID, menuCode, &btnList)
		if err != nil {
			response.FailWithMessage("获取操作列表失败", c)
			return
		}
	}
	response.OkWithDetailed(&btnList, "获取操作列表成功", c)
}

// @Tags 菜单
// @Summary 获取所有菜单
// @Produce  application/json
// @Security ApiKeyAuth
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取用户信息成功"}"
// @Router /api/v1/role/allrole [get]
func (h *handler) AllMenu(c *gin.Context) {
	menuList,err := h.dao.GetAllMenu()
	if err != nil {
		log.Errorf("获取失败,%v",err)
		response.FailWithMessage("获取失败",c)
		return
	}
	resData := menuList
	response.OkWithDetailed(resData, "添加成功", c)
}
