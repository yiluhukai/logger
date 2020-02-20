/*
 	1.可以全局设置日志等级,只有高于日志等级的日志才会打印
	2.可以根据配置文件选择打印日志的位置
	3.对日志库进行封装，方便对外调用
	4.对输入到文件的日志进行异步打印
	5.对日志文件进行切割，可以按小时切割和文件的打印切割
 */

package logger

import (
	"fmt"
	"time"
)

var log Logger
//定义日志接口:可以使用的文件日志和console日志
type Logger interface {
	SetLevel(level int)
	LogDebug(format string,args ...interface{})
	LogTrace(format string,args ...interface{})
	LogInfo(format string,args ...interface{})
	LogWarn(format string,args ...interface{})
	LogError(format string,args ...interface{})
	//程序不能运行级别的
	LogFatal(format string,args ...interface{})
	Close()
}


//初始化日志

func InitLogger(name string,pathName string,fileName string,level int,log_split_type string) (err error){
	//使用map可以底层对config的不一致
	config := make(map[string]string,8)
	config["path_name"] = pathName
	config["file_name"] =fileName
	config["level"] = getLoggerLevelText(level)
	config["log_split_type"] = log_split_type
	switch name {
	case "file":
		log,err = NewFileLogger(config)
	case "console":
		log,err = NewConsoleLogger(config)
	default:
		err=fmt.Errorf("unkown logger")
	}
	return
}

func Run(){
	for{
		log.LogDebug("%s","hello world")
		time.Sleep(time.Second*60)
	}
}

func LogDebug(format string,args ...interface{}){
	log.LogDebug(format,args...)
}

func LogTrace(format string,args ...interface{}){
	log.LogTrace(format,args...)
}

func LogInfo(format string,args ...interface{}){
	log.LogInfo(format,args...)
}

func LogWarn(format string,args ...interface{}){
	log.LogWarn(format,args...)
}

func LogError(format string,args ...interface{}){
	log.LogError(format,args...)
}

func LogFatal(format string,args ...interface{}){
	log.LogFatal(format,args...)
}

