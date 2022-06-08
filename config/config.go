package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Runmode string
	Addr    string
}

type JwtConfig struct {
	JwtSecret string `json:"jwt_secret"`
}

type MysqlConfig struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int `json:"max_idle_connections"`
	MaxOpenConnections    int `json:"max_open_connections"`
	MaxConnectionLifeTime int `json:"max_connection_life_time"`
	LogLevel              int `json:"log_level"`
}

type Config struct {
	ServerConfig ServerConfig
	JwtConfig    JwtConfig
	MysqlConfig  MysqlConfig
}

var Cfg Config

func init() {
	viper.SetConfigName("douyin")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("D:/go/go-douyin-demo/config/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config file not found")
		} else {
			panic("read config error: " + err.Error())
		}
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic("unmarshal douyin.yaml error:" + err.Error())
	}
}
