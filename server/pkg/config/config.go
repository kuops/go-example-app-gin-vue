package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Redis RedisConfig `mapstructure:"redis"`
	Mysql MySQLConfig `mapstructure:"mysql"`
	Kubernetes KubernetesConfig `mapstructure:"kubernetes"`
}

func InitializeConfig(cfgFile string) *Config {
	var config = &Config{}
	var err error

	if cfgFile == "" {
		cfgFile = "configs/dev.yaml"
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgFile)

	if err = viper.ReadInConfig();err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	return config
}