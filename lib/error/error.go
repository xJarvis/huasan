package error

import (
	"fmt"
	"github.com/xjarvis/huashan/log/logger"
	"runtime"
)


func Catch() {
	if err:=recover(); err!=nil {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			logger.Info(fmt.Sprintf("#%s#%s#%d line", file, runtime.FuncForPC(pc).Name(), line))
		}
		logger.Error("catch error:",err)
	}
}