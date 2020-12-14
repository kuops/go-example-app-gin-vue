package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
)

const groupName = "user"

func Register(group *gin.RouterGroup,mysqlClient *mysql.Client) {
	rg := group.Group(groupName)
	rg.POST("/login", LoginHandlerFunc(mysqlClient))
}