package store

import (
	"go-douyin-demo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var Db = Init()

//GetLogLever 获取日志级别
func GetLogLever(logLevel int) logger.LogLevel {
	switch logLevel {
	case 1:
		return logger.Silent
	case 2:
		return logger.Error
	case 3:
		return logger.Warn
	default:
		return logger.Info
	}
}

func Init() *gorm.DB {
	var mysqlCfg = config.Cfg.MysqlConfig
	dsn := mysqlCfg.Username + ":" + mysqlCfg.Password + "@tcp(" + mysqlCfg.Host + ")/" + mysqlCfg.Database +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(GetLogLever(mysqlCfg.LogLevel)),
	})
	if err != nil {
		panic("db connect error :" + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConnections)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConnections)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxIdleTime(time.Second * 10)

	Migrate(db)

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Video{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&FavoriteVideo{})
	db.AutoMigrate(&UserRelation{})
}
