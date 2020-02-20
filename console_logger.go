package logger

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(config map[string]string) (log Logger,err error){
	level,ok:=config["level"]
	if !ok{
		err=fmt.Errorf("level is invalid")
		return
	}
	log = &ConsoleLogger{
		level:getLoggerLevel(level),
	}
	return
}


func (logger * ConsoleLogger) SetLevel(level int)  {
	if level<LogLevelDebug || level>LogLevelFatal{
		logger.level = LogLevelDebug
		return
	}
	logger.level = level
}

func (logger * ConsoleLogger) LogDebug(format string,args ...interface{})  {

	chanData:=printfLogger(logger.level,LogLevelDebug,format,args...)
	if chanData ==nil{
		return
	}else{
		log:=fmt.Sprintf("%s %s (%s %s %d) %s",chanData.timeStr,
			chanData.levelStr,chanData.fileName,chanData.funcName,chanData.lineNum,chanData.message)
		fmt.Fprintln(os.Stdout,log)
	}
}

func (logger * ConsoleLogger) LogTrace(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelTrace,format,args...)
	if chanData ==nil{
		return
	}else{
		log:=fmt.Sprintf("%s %s (%s %s %d) %s",chanData.timeStr,
			chanData.levelStr,chanData.fileName,chanData.funcName,chanData.lineNum,chanData.message)
		fmt.Fprintln(os.Stdout,log)
	}
}

func (logger * ConsoleLogger) LogInfo(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelInfo,format,args...)
	if chanData ==nil{
		return
	}else{
		log:=fmt.Sprintf("%s %s (%s %s %d) %s",chanData.timeStr,
			chanData.levelStr,chanData.fileName,chanData.funcName,chanData.lineNum,chanData.message)
		fmt.Fprintln(os.Stdout,log)
	}
}

func (logger * ConsoleLogger) LogWarn(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelWarn,format,args...)
	if chanData ==nil{
		return
	}else{
		log:=fmt.Sprintf("%s %s (%s %s %d) %s",chanData.timeStr,
			chanData.levelStr,chanData.fileName,chanData.funcName,chanData.lineNum,chanData.message)
		fmt.Fprintln(os.Stdout,log)
	}
}

func (logger * ConsoleLogger) LogError(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelError,format,args...)
	if chanData ==nil{
		return
	}else{
		log:=fmt.Sprintf("%s %s (%s %s %d) %s",chanData.timeStr,
			chanData.levelStr,chanData.fileName,chanData.funcName,chanData.lineNum,chanData.message)
		fmt.Fprintln(os.Stdout,log)
	}
}

func (logger * ConsoleLogger) LogFatal(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelFatal,format,args...)
	if chanData ==nil{
		return
	}else{
		log:=fmt.Sprintf("%s %s (%s %s %d) %s",chanData.timeStr,
			chanData.levelStr,chanData.fileName,chanData.funcName,chanData.lineNum,chanData.message)
		fmt.Fprintln(os.Stdout,log)
	}
}

func (logger * ConsoleLogger) Close(){

}





