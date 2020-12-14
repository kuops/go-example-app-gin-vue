package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/middleware/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Register(group *gin.RouterGroup) {
	group.GET("/metrics",prometheus.PromHandler(promhttp.Handler()))
}