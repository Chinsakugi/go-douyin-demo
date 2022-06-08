package test

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func TestGetViperConfig(t *testing.T) {
	viper.SetConfigFile("../config/douyin.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(viper.Get("addr"))
}

func TestViperConfig(t *testing.T) {
	fmt.Println(config.Cfg)
}
