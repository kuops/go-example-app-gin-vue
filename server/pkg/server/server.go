package server

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/kuops/go-example-app/server/docs"
	"github.com/kuops/go-example-app/server/pkg/apis/metrics"
	userv1 "github.com/kuops/go-example-app/server/pkg/apis/user/v1"
	"github.com/kuops/go-example-app/server/pkg/config"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/middleware/cors"
	"github.com/kuops/go-example-app/server/pkg/middleware/prometheus"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/utils/ip"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Server struct {
	Server       *http.Server
	ServerConfig *config.ServerConfig
	MySQLClient  *mysql.Client
}

func (s *Server) PrepareRun() error {
	s.setMode()
	s.installRouters()
	log.Infof("注册路由成功...")
	return nil
}

func (s *Server) Run(stopCh <-chan struct{}) error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = s.Server.ListenAndServe()
	}()

	if err != nil {
		return err
	}

	log.Infof("服务运行中,访问地址: http://%v:%v", ip.GetIPAddress(), s.ServerConfig.Port)
	log.Infof("swagger 页面: http://%v:%v/swagger/index.html", ip.GetIPAddress(), s.ServerConfig.Port)
	log.Infof("metrics 页面: http://%v:%v/metrics", ip.GetIPAddress(), s.ServerConfig.Port)

	select {
	case <-stopCh:
		log.Info("正在关闭服务....")
		_ = s.Server.Shutdown(ctx)
	}

	return nil
}

func (s *Server) installRouters() {
	router := gin.New()

	router.Use(cors.Middleware())
	router.Use(prometheus.PromMiddleware(&prometheus.PromOpts{}))

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	rootGroup := router.Group("")
	v1Group := router.Group("api/v1")

	metrics.Register(rootGroup)
	userv1.Register(v1Group, s.MySQLClient)

	s.Server.Handler = router
}

func (s *Server)setMode()  {
	if s.ServerConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}