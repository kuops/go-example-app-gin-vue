package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/request"
	"github.com/kuops/go-example-app/server/pkg/response"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"github.com/kuops/go-example-app/server/pkg/utils/jwt"
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
// @Param data body request.Login true "用户登录接口"
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
			ID: u.ID,
			UUID: key,
			Name: user.UserName,
		}
		cliams.ExpiresAt = expires
		token,err := jwt.CreateToken(cliams)
		if err != nil {
			log.Errorf("创建 Token 失败",err)
			response.FailWithMessage("服务内部错误", c)
			return
		}

		data := map[string]interface{}{
			"user": user,
			"token": token,
			"exp": expires,
		}

		response.OkWithDetailed(data,"登录成功", c)
}

// @Tags 用户
// @Summary 用户退出
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"退出成功"}"
// @Router /api/v1/user/logout [get]
func (h *handler)Logout(c *gin.Context) {
	cliamsContext,_ := c.Get("claims")
	claims := cliamsContext.(*jwt.CustomClaims)
	_ = h.cache.Del(claims.UUID)
	response.OkWithMessage("退出成功",c)
}