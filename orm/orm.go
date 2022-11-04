package orm

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"huashan/config"
	llog "huashan/logger"
	"strings"
	"time"
)

const (
	CHARACTER 		= "utf8"

	MaxIdleConn 	= 128
	MaxOpenConn 	= 512
	MaxLifeTime     = 3600
)

var (
	DB *gorm.DB
)

// 获取数据库引擎DSN
func getDbEngineDSN(sEngine string) string {
	dsn := ""
	switch sEngine {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			config.Get("database.username"),
			config.Get("database.password"),
			config.Get("database.host"),
			config.Get("database.database"),
			CHARACTER)
	default:
		panic(errors.New("database engine not exists! "))
	}
	return dsn
}


func NewOrm() *gorm.DB {
	if DB != nil {
		return DB
	}

	gLogger := logger.New(llog.Logger,logger.Config{
		SlowThreshold: 100 * time.Millisecond,
		LogLevel:      logger.Info,
		Colorful:      true,
	})

	sEngine := strings.ToLower(config.Get("database.engine"))
	db, err := gorm.Open(mysql.Open(getDbEngineDSN(sEngine)), &gorm.Config{
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true, //禁用外键约束
		Logger:                                   gLogger,
	})

	if err != nil {
		panic(err)
	}

	sqlDB,err := db.DB()
	if  err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(MaxIdleConn)
	sqlDB.SetMaxIdleConns(MaxOpenConn)
	sqlDB.SetConnMaxLifetime(MaxLifeTime * time.Second)

	return db
}

func Initialize() {
	DB = NewOrm()
	llog.Info("gorm init finish")
}

func UnInitialize() {
	if DB != nil {
		DB = nil
	}
	llog.Info("gorm release finish")
}