package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
)

func LoginHandlerFunc(mysqlClient *mysql.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbName := mysqlClient.Database().DB.Name()
		c.JSON(200,gin.H{"database": dbName})
	}
}
