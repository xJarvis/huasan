package orm

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"huashan/config"
	"huashan/logger"
	"strings"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	MAXIDLECONN 	= 600
	MAXOPENCONN 	= 600
	CHARACTER 		= "utf8"
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
	sEngine := strings.ToLower(config.Get("database.engine"))
	db, err := gorm.Open(sEngine, getDbEngineDSN(sEngine))
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(MAXIDLECONN)
	db.DB().SetMaxIdleConns(MAXOPENCONN)
	db.SetLogger(&logger.Logger)

	// 开发环境开启sql日志
	if config.DEBUG == config.Get("command.env") {
		db.LogMode(true)
	}

	return db
}

func Initialize() {
	DB = NewOrm()
	logger.Info("gorm init finish")
}

func UnInitialize() {
	if DB != nil {
		DB.Close()
	}
	logger.Info("gorm release finish")
}