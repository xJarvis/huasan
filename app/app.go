package app

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xjarvis/huashan/app/command"
	config2 "github.com/xjarvis/huashan/app/config"
	error2 "github.com/xjarvis/huashan/lib/error"
	"github.com/xjarvis/huashan/lib/signal"
    "github.com/xjarvis/huashan/lib/file/directory"
	"github.com/xjarvis/huashan/log/logger"
	"os"
	"path/filepath"
	"syscall"
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

	logger.Info("app env release finish!")
	fmt.Println("app env release finish!")
}

// SignalRun /**退出信号捕获处理*/
func SignalRun() {
	signal.RegisterSignal(syscall.SIGHUP,ReleaseEnv)
	signal.RegisterSignal(syscall.SIGINT,ReleaseEnv)
	signal.RegisterSignal(syscall.SIGQUIT,ReleaseEnv)
	signal.CatchSignal()
}