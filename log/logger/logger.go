package logger

import (
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"runtime"
)

// 日志库



type Level int8

var debug bool		//调试模式标识
var logger seelog.LoggerInterface

const (
	DEBUG = iota 	//调试
	INFO			//普通
	WARN			//警告
	ERROR			//错误
	FATAL			//致命
)

type Model struct {
	Content string `json:"content"`			//前缀
}

func (p *Model)Debug(v ...interface{}) {
	write(DEBUG, p.Content,v...)
}

func (p *Model)Print(v ...interface{}) {
	write(INFO, p.Content,v...)
}

func (p Model)Printf(s string,v ...interface{}) {
	write(INFO, p.Content,fmt.Sprintf(s,v...))
}

func (p *Model)Info(v ...interface{}) {
	write(INFO, p.Content,v...)
}

func (p *Model)Infof(s string,v ...interface{}) {
	write(INFO, p.Content,fmt.Sprintf(s,v...))
}

func (p *Model)Warn(v ...interface{}) {
	write(WARN, p.Content,v...)
}

func (p *Model)Warnf(s string,v ...interface{}) {
	write(WARN, p.Content,fmt.Sprintf(s,v...))
}

func (p *Model)Error(v ...interface{}) {
	write(ERROR, p.Content,v...)
}

func (p *Model)Errorf(s string,v ...interface{}) {
	write(ERROR, p.Content,fmt.Sprintf(s,v...))
}

func (p *Model)Fatal(v ...interface{}) {
	write(FATAL, p.Content,v...)
}

var Logger Model

func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

func Print(v ...interface{}) {
	Logger.Print(v...)
}

func Printf(s string, v ...interface{}) {
	Logger.Printf(s,v...)
}

func Info(v ...interface{}) {
	Logger.Info(v...)
}

func Infof(s string,v ...interface{}) {
	Logger.Infof(s, v...)
}

func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

func Warnf(s string,v ...interface{}) {
	Logger.Warnf(s, v...)
}

func Error(v ...interface{}) {
	Logger.Error(v...)
}

func Errorf(s string,v ...interface{}) {
	Logger.Errorf(s, v...)
}

func Fatal(v ...interface{}) {
	Logger.Fatal(v...)
}

func Initialize(logFile string, isDebug bool) {
	debug = isDebug
	var err error
	logger, err = seelog.LoggerFromConfigAsString(getCfgXml(logFile))
	if err != nil {
		panic(err)
	}

	Logger = Model{""}
}

func UnInitialize() {
	if logger != nil {
		logger.Close()
	}
}

func write(level Level,content string, v ...interface{}) {
	defer logger.Flush()
	switch level {
	case DEBUG:
		logger.Debug(content, v)
	case INFO:
		logger.Info(content, v)
	case WARN:
		logger.Warn(content, v)
	case FATAL:
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			content += fmt.Sprintf("#%s#%s#%d行#", file, runtime.FuncForPC(pc).Name(), line)
		}
		logger.Critical(content, v)
		os.Exit(1)//致命错误直接退出
	case ERROR:
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			content += fmt.Sprintf("#%s#%s#%d行#", file, runtime.FuncForPC(pc).Name(), line)
		}
		logger.Error(content, v)
	}
}


func getCfgXml(logFile string) string {
	cfgXml := `
   <seelog type="asynctimer" asyncinterval="5000000" minlevel="debug" maxlevel="error">
       <outputs formatid="main">
			%s
           <filter levels="info,critical,warn,error">
               <rollingfile formatid="main" type="size" filename="%s" maxsize="100000000" maxrolls="100" />
           </filter>
       </outputs>
       <formats>
           <format id="main" format="%%Date %%Time [%%LEV] %%Msg%%n"/>
       </formats>
   </seelog>
   `
	filterStr := ""
	if debug {
		filterStr = `
			<filter levels="info,debug,critical,warn,error">
               <console />
           	</filter>
			`
	}
	cfgXml = fmt.Sprintf(cfgXml,filterStr,logFile)
	return cfgXml
}




