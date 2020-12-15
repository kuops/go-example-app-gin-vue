package cmd

import (
	"fmt"
	"github.com/kuops/go-example-app/server/pkg/config"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/server"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	"github.com/spf13/cobra"
	"net/http"
)

var cfgFile string

func init()  {
	log.InitLogger()
}

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "app",
		Long: `the simple golang restful web framework demo project`,
		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cfgFile,server.SetupSignalHandler())
		},
	}

	fs := cmd.Flags()
	fs.StringVar(&cfgFile,"config","","config file path. (default: is configs/dev.yaml)")

	return cmd
}

func NewServer(config *config.Config,stopCh <-chan struct{}) (*server.Server, error) {

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d",config.Server.Port),
	}

	s := &server.Server{
		ServerConfig: &config.Server,
		Server: httpServer,
	}

	mysqlClient,err := mysql.NewMySQLClient(&config.Mysql,stopCh)
	if err != nil {
		return nil, err
	}
	s.MySQLClient = mysqlClient

	redisClient,err := redis.NewRedisClient(&config.Redis,stopCh)
	if err != nil {
		return nil, err
	}
	s.RedisClient = redisClient

	return s, nil
}

func Run(cfg string,stopCh <-chan struct{}) error {
	appConfig := config.InitializeConfig(cfg)
	log.Info("初始化配置文件成功...")

	apiServer, err := NewServer(appConfig,stopCh)
	if err != nil {
		return err
	}

	err = apiServer.PrepareRun()
	if err != nil {
		return err
	}

	err = apiServer.Run(stopCh)
	if err != nil {
		return err
	}
	return nil
}


