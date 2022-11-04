package signal

import (
	"github.com/xjarvis/huashan/log/logger"
	"os"
	"os/signal"
	"syscall"
)

// 信号初始化
var (
	sysChan  		chan os.Signal
	funcSighup 		func()
	funcSigint 		func()
	funcSigterm  	func()
)

// 初始化通道
func init() {
	sysChan = make(chan os.Signal)
}
// 阻塞运行
func BolockRuning() {
	signal.Notify(sysChan, syscall.SIGINT)
	<-sysChan
}

// 注册信号
func RegisterSignal(sig os.Signal,sigFunc func()) {
	switch sig {
	case syscall.SIGHUP:
		funcSighup = sigFunc
	case syscall.SIGINT:
		funcSigint = sigFunc
	case syscall.SIGQUIT:
		funcSigterm = sigFunc
	}
}


// 捕捉信号
func CatchSignal() {
	signal.Notify(sysChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	for {
		s := <-sysChan
		switch s {
		case syscall.SIGHUP:
			logger.Info("signal:terminal disconnect")
			if funcSighup != nil {
				funcSighup()}
		case syscall.SIGINT:
			logger.Info("signal:ctrl+c")
			if funcSigint != nil {
				funcSigint()}
		case syscall.SIGQUIT:
			logger.Info("signal:application quit")
			if funcSigterm != nil {
				funcSigterm()}
		}
	}
}