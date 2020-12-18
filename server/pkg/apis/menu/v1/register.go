package v1

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
)

const groupName = "menu"

func Register(group *gin.RouterGroup,mysqlClient *mysql.Client,redisClient redis.Interface,enforcer *casbin.Enforcer) {
	handler := newHandler(mysqlClient,redisClient,enforcer)

	rg := group.Group(groupName)
	rg.POST("/list",handler.List)
	rg.GET("/detail",handler.Detail)
	rg.POST("/update",handler.Update)
	rg.POST("/create",handler.Create)
	rg.GET("/menubuttonlist",handler.MenuButtonList)
	rg.POST("/delete",handler.Delete)
	rg.GET("/allmenu",handler.AllMenu)
}

