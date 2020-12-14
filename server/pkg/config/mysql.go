package config

type MySQLConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

