package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
)

const groupName = "user"

func Register(group *gin.RouterGroup,mysqlClient *mysql.Client,redisClient redis.Interface) {
	handler := newHandler(mysqlClient,redisClient)

	rg := group.Group(groupName)
	rg.POST("/login", handler.Login)
	rg.POST("/logout",handler.Logout)
	rg.GET("/info",handler.Info)
}

