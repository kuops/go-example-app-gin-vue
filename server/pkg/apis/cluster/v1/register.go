package v1

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

const groupName = "cluster"

func Register(group *gin.RouterGroup,mysqlClient *mysql.Client,redisClient redis.Interface,enforcer *casbin.Enforcer,clientset *kubernetes.Clientset,dynamic dynamic.Interface) {
	handler := newHandler(mysqlClient,redisClient,enforcer,clientset,dynamic)
	rg := group.Group(groupName)
	rg.POST("/namespace/list",handler.ListNamespace)
	rg.POST("/namespace/create",handler.CreateNamespace)
	rg.POST("/namespace/delete",handler.DeleteNamespace)
	rg.POST("/deployment/list",handler.ListDeployment)
	rg.POST("/deployment/update",handler.UpdateDeployment)
	rg.POST("/deployment/delete",handler.DeleteDeployment)
	rg.GET("/deployment/create",handler.CreateDeployment)
}
