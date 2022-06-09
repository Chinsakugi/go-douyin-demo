package test

import (
	"fmt"
	"go-douyin-demo/config"
	"go-douyin-demo/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGorm(t *testing.T) {
	var mysqlCfg = config.Cfg.MysqlConfig
	dsn := mysqlCfg.Username + ":" + mysqlCfg.Password + "@tcp(" + mysqlCfg.Host + ")/" + mysqlCfg.Database +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("db connect error :" + err.Error())
	}
	fmt.Println(db)
}

func TestCreateTable(t *testing.T) {
	table := store.Db.Migrator().HasTable("user")
	fmt.Println(table)
}
