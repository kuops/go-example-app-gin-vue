package config

type ServerConfig struct {
	Port int `mapstructure:"port"`
	Debug bool `mapstructure:"debug"`
}