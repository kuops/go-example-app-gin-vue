package v1

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
)

const groupName = "user"

func Register(group *gin.RouterGroup,mysqlClient *mysql.Client,redisClient redis.Interface,enforcer *casbin.Enforcer) {
	handler := newHandler(mysqlClient,redisClient,enforcer)

	rg := group.Group(groupName)
	rg.POST("/login", handler.Login)
	rg.POST("/logout",handler.Logout)
	rg.GET("/info",handler.Info)
	rg.POST("/changePassword",handler.ChangePassword)
	rg.POST("/list",handler.List)
	rg.GET("/detail",handler.Detail)
	rg.POST("/update",handler.Update)
	rg.POST("/create",handler.Create)
	rg.POST("/delete",handler.Delete)
	rg.GET("/usersroleidlist",handler.UsersRoleList)
	rg.POST("/setrole",handler.SetRole)
}

