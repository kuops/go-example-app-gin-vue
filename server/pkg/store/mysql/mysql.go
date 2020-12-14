package mysql

import (
	"fmt"
	"github.com/kuops/go-example-app/server/pkg/config"
	"github.com/kuops/go-example-app/server/pkg/log"
	"gorm.io/driver/mysql"
	"time"
)
import "gorm.io/gorm"

type Database struct {
	DB *gorm.DB
}

type Client struct {
	database *Database
}

func NewMySQLClient(cfg *config.MySQLConfig,stopCh <-chan struct{}) (*Client, error) {
	var client Client

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/demo?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port)
	mysqlConfig := mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{}); err != nil {
		log.Errorf("database connection error: %v\n",err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(2)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(time.Second * 300)

		client.database = &Database{
			DB: db,
		}

		go func() {
			<-stopCh
			log.Info("断开数据库连接....")
			if err := sqlDB.Close(); err != nil {
				log.Warnf("error happened during closing mysql connection: %v\n", err)
			}
		}()

		log.Infof("数据库初始化连接成功，连接信息: %v:%v",cfg.Host,cfg.Port)
	}

	return &client,nil
}

func (client *Client) Database() *Database {
	return client.database
}