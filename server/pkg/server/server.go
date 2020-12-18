package server

import (
	"context"
	"github.com/gin-gonic/gin"
	menuv1 "github.com/kuops/go-example-app/server/pkg/apis/menu/v1"
	"github.com/kuops/go-example-app/server/pkg/apis/metrics"
	rolev1 "github.com/kuops/go-example-app/server/pkg/apis/role/v1"
	userv1 "github.com/kuops/go-example-app/server/pkg/apis/user/v1"
	"github.com/kuops/go-example-app/server/pkg/casbin"
	"github.com/kuops/go-example-app/server/pkg/config"
	dbinit "github.com/kuops/go-example-app/server/pkg/db/initialize"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/middleware/auth"
	"github.com/kuops/go-example-app/server/pkg/middleware/cors"
	"github.com/kuops/go-example-app/server/pkg/middleware/logger"
	"github.com/kuops/go-example-app/server/pkg/middleware/prometheus"
	casbinmw "github.com/kuops/go-example-app/server/pkg/middleware/casbin"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"github.com/kuops/go-example-app/server/pkg/utils/ip"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Server struct {
	Server       *http.Server
	ServerConfig *config.ServerConfig
	MySQLClient  *mysql.Client
	MySQLConfig  *config.MySQLConfig
	RedisClient  redis.Interface
	Casbin     *casbin.Casbin
}

func (s *Server) PrepareRun() error {
	s.setMode()
	s.Migration()
	s.InitialDatas()
	s.InitCsbinEnforcer()
	s.installRouters()


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

	skipAuthUri := []string{
		"/api/v1/user/login",
		"/swagger",
		"/metrics",
	}

	skipCasbinUri := []string {
		"/api/v1/user/login",
		"/swagger",
		"/metrics",
		"/api/v1/user/info",
		"/api/v1/user/roleList",
		"/api/v1/user/changePassword",
		"/api/v1/menu/allmenu",
		"/api/v1/menu/menubuttonlist",
		"/api/v1/role/allrole",
		"/api/v1/role/rolemenuidlist",
	}

	middlewares := []gin.HandlerFunc {
		logger.GinLogger(),
		logger.GinRecovery(),
		cors.Middleware(),
		auth.Middleware(s.RedisClient,skipAuthUri),
		prometheus.PromMiddleware(&prometheus.PromOpts{}),
		casbinmw.Middleware(skipCasbinUri,s.Casbin.Enforcer),
	}

	router.Use(middlewares...)


	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	rootGroup := router.Group("")
	apiv1Group := router.Group("api/v1")

	metrics.Register(rootGroup)
	userv1.Register(apiv1Group, s.MySQLClient,s.RedisClient,s.Casbin.Enforcer)
	menuv1.Register(apiv1Group, s.MySQLClient,s.RedisClient,s.Casbin.Enforcer)
	rolev1.Register(apiv1Group, s.MySQLClient,s.RedisClient,s.Casbin.Enforcer)

	s.Server.Handler = router
	log.Info("注册路由成功...")
}

func (s *Server)setMode()  {
	if s.ServerConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func (s *Server)Migration() {
	log.Info("初始化数据库表结构....")
	dbinit.Migration(s.MySQLClient)
}

func (s *Server)InitialDatas()  {
	log.Info("初始化数据库系统数据...")
	dbinit.InitialDatas(s.MySQLClient)
}

func (s *Server)InitCsbinEnforcer()  {
	log.Info("初始化权限系统...")
	s.Casbin.InitCsbinEnforcer(s.MySQLConfig,s.MySQLClient)
}