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
	"gorm.io/gorm"
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

// @Tags 角色
// @Summary 分页角色列表
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/v1/role/list [post]
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
	list,total,err :=  h.dao.GetRolesList(limit,offset)

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

// @Tags 角色
// @Summary 角色详情
// @Produce  application/json
// @Security ApiKeyAuth
// @Param id query int false "int valid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取角色详情成功"}"
// @Router /api/v1/role/detail [get]
func (h *handler)Detail(c *gin.Context) {
	var id uint64
	var err error
	var role Role
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
	role.ID = id
	m,err := h.dao.GetRolesDetail(&role)
	if err != nil {
		log.Errorf("获取失败 %v",err)
		response.FailWithMessage("获取失败", c)
	}

	resData :=  m
	response.OkWithDetailed(resData, "获取成功", c)
}


// @Tags 角色
// @Summary 更新角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body Role true "ID, 父级ID, URL"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新角色成功"}"
// @Router /api/v1/role/update [post]
func (h *handler) Update(c *gin.Context) {
	role := Role{}
	_ = c.ShouldBindJSON(&role)
	if role.ID == 0 {
		response.FailWithMessage("菜单 ID 不能为空", c)
		return
	}

	updatedRole,err := h.dao.UpdateRoleInfo(&role)
	if err != nil {
		log.Errorf("更新失败 %v",err)
		response.FailWithMessage("更新失败", c)
		return
	}

	resData := updatedRole
	response.OkWithDetailed(resData, "更新成功", c)
}


// @Tags 角色
// @Summary 添加角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body Menu true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /api/v1/role/create [post]
func (h *handler) Create(c *gin.Context) {
	req := Role{}
	err := c.ShouldBindJSON(&req)

	if err != nil {
		response.FailWithMessage("缺少必填项", c)
		return
	}

	m,err := h.dao.CreateRole(&req)

	if err != nil {
		log.Errorf("添加角色失败, %v",err)
		response.FailWithMessage("添加角色失败", c)
		return
	}

	resData := m
	response.OkWithDetailed(resData, "添加成功", c)
}

// @Tags 角色
// @Summary 删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteRoles true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除角色成功"}"
// @Router /api/v1/role/delete [post]
func (h *handler) Delete(c *gin.Context) {
	req := request.DeleteRoles{}
	var err error
	err = c.ShouldBindJSON(&req)
	if err != nil || len(req.IDS) == 0 {
		log.Errorf("删除失败,%v",err)
		response.FailWithMessage("删除失败",c)
		return
	}

	if err := h.dao.DeleteRoles(req.IDS); err != nil {
		log.Errorf("删除失败,%v",err)
		response.FailWithMessage("删除失败",c)
		return
	}
	go CsbinDeleteRole(h.enforcer,req.IDS)
	response.OkWithMessage("删除用户成功", c)
}

func CsbinDeleteRole(enforcer *casbin.Enforcer,roleids []uint64) {
	for _, rid := range roleids {
		enforcer.DeletePermissionsForUser(convert.ToString(rid))
		enforcer.DeleteRole(convert.ToString(rid))
	}
	enforcer.SavePolicy()
}

// @Tags 角色
// @Summary 获取所有角色
// @Produce  application/json
// @Security ApiKeyAuth
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取用户信息成功"}"
// @Router /api/v1/role/allrole [get]
func (h *handler) AllRole(c *gin.Context) {
	roleList,err := h.dao.GetAllRole()
	if err != nil {
		log.Errorf("获取失败,%v",err)
		response.FailWithMessage("获取失败",c)
		return
	}
	resData := roleList
	response.OkWithDetailed(resData, "添加成功", c)
}

// @Tags 角色
// @Summary 获取角色下拉菜单
// @Produce  application/json
// @Security ApiKeyAuth
// @Param roleid query int false "int valid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取用户信息成功"}"
// @Router /api/v1/role/rolemenuidlist [get]
func (h *handler) RoleMenuIDList(c *gin.Context) {
	id,_ := c.GetQuery("roleid")
	roleID := convert.ToUint64(id)
	menuIDList := []uint64{}
	where := RoleMenu{RoleID: roleID}
	err := h.dao.GetRoleMenuIDList(&RoleMenu{}, &where, &menuIDList, "menu_id")
	if err != nil {
		log.Errorf("获取失败,%v",err)
		response.FailWithMessage("获取失败",c)
		return
	}
	resData := menuIDList
	response.OkWithDetailed(resData, "获取成功", c)
}

// @Tags 角色
// @Summary 设置角色菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SetRoleMenus true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除用户成功"}"
// @Router /api/v1/role/setrole [post]
func (h *handler) SetRoleMenus(c *gin.Context) {
	req := request.SetRoleMenus{}
	err := c.Bind(&req)
	if err != nil {
		log.Errorf("缺少必填字段")
		response.FailWithMessage("缺少必填字段",c)
		return
	}
	err = h.dao.SetRoleMenus(req.RoleID,req.MenuIDS)
	if err != nil {
		log.Errorf("设置角色菜单失败,%v",err)
		response.FailWithMessage("设置角色菜单失败",c)
		return
	}
	go CasbinSetRolePermission(h.dao.db,h.enforcer,req.RoleID)

	response.OkWithMessage("设置成功", c)
}

func CasbinSetRolePermission(db *gorm.DB,enforcer *casbin.Enforcer,roleid uint64) {
	var roleMenus []RoleMenu
	err := db.Where(&RoleMenu{RoleID: roleid}).Find(&roleMenus).Error
	if err != nil {
		return
	}
	for _, roleMenu := range roleMenus {
		menu := Menu{}
		where := Menu{}
		where.ID = roleMenu.MenuID
		err = db.Where(&where).First(&menu).Error
		if err != nil {
			return
		}
		if  menu.MenuType == 3 {
			_,err := enforcer.AddPermissionForUser(convert.ToString(roleid), "/api/v1" + menu.URL,"GET|POST|UPDATE|DELETE")
			if err != nil {
				log.Error(err)
			}
		}
	}
}