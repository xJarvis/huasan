package app

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xjarvis/huashan/app/command"
	config2 "github.com/xjarvis/huashan/app/config"
	"github.com/xjarvis/huashan/cache/local"
	"github.com/xjarvis/huashan/cache/redis"
	error2 "github.com/xjarvis/huashan/lib/error"
	"github.com/xjarvis/huashan/lib/exsignal"
	"github.com/xjarvis/huashan/lib/file/directory"
	"github.com/xjarvis/huashan/log/logger"
	"github.com/xjarvis/huashan/mysql/orm"
	"github.com/xjarvis/huashan/nosql/influxdb"
	"github.com/xjarvis/huashan/nosql/mongo"
	"github.com/xjarvis/huashan/schedule/cron"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	AppDir 	string // 应用根目录
)

// InitEnv /**初始化程序运行环境*/
func InitEnv() {
	var err error

	//获取执行目录
	AppDir, err = directory.GetWorkDir()
	if err != nil {
		panic(err)
	}

	// 命令行参数获取
	mCmd := command.Command()
	if _,ok := mCmd["cfg"]; !ok {
		panic(errors.New("cfg file not exit!"))
	}

	//配置文件导入
	if !filepath.IsAbs(mCmd["cfg"]) {
		mCmd["cfg"] = filepath.Join(AppDir, mCmd["cfg"])
	}
	config2.InitConfig(mCmd["cfg"])
	for k,cmd := range mCmd {
		config2.SetValue("command" + "." +  k,cmd)
	}

	if "true" == config2.Get("system.debug") {
		config2.SetValue("command.env","dev")
	}

	//日志文件导入
	logFile,err := config2.GetValue("log.file")
	if err != nil {
		panic(err)
	}
	if !filepath.IsAbs(logFile) {
		logFile = filepath.Join(AppDir, logFile)
	}

	//日志文件初始化
	if config2.DEBUG == config2.Get("command.env") {
		logger.Initialize(logFile,true)
	} else {
		logger.Initialize(logFile,false)
	}

	//缓存启动
	if config2.Get("cache.open") == "yes" {
		local.Initialize(time.Duration(config2.GetInt("cache.expire")) * time.Second)
	}

	//Influx 启动
	if config2.Get("influx.open") == "yes" {
		influxdb.Initialize()
		influxdb.Run()
	}

	// Redis 初始化
	if config2.Get("redis.open") == "yes" {
		redis.Initialize()
	}
	//DB ORM初始化
	if config2.Get("database.open") == "yes" {
		orm.Initialize()
	}

	//mongo 初始化
	if config2.Get("mgo.open") == "yes" {
		mongo.Initialize()
	}

	//计划任务
	if config2.Get("task.open") == "yes" {
		cron.Initialize()
	}

	/*todo 核心业务环境启动*/

	logger.Info("app env init finish!")
}

// ReleaseEnv /**释放运行环境*/
func ReleaseEnv() {
	//logger释放
	defer func() {
		error2.Catch()
		logger.UnInitialize()
		fmt.Println("application exit!")
		os.Exit(0) // 程序退出
	}()

	/*todo 核心业务环境释放*/

	//计划任务释放
	if config2.Get("task.open") == "yes" {
		cron.UnInitialize()
	}

	//MONGO 釋放
	if config2.Get("mgo.open") == "yes" {
		mongo.UnInitialize()
	}

	//DB ORM释放
	if config2.Get("database.open") == "yes" {
		orm.UnInitialize()
	}

	// Redis释放
	if config2.Get("redis.open") == "yes" {
		redis.UnInitialize()
	}

	//缓存釋放
	if config2.Get("cache.open") == "yes" {
		local.UnInitialize()
	}

	logger.Info("app env release finish!")
	fmt.Println("app env release finish!")
}

// SignalRun /**退出信号捕获处理*/
func SignalRun() {
	exsignal.RegisterSignal(syscall.SIGHUP,ReleaseEnv)
	exsignal.RegisterSignal(syscall.SIGINT,ReleaseEnv)
	exsignal.RegisterSignal(syscall.SIGQUIT,ReleaseEnv)
	exsignal.CatchSignal()
}