package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

//定义一个用来存放管道中信息的结构体
type ChanData struct{
	message string
	timeStr string
	levelStr string
	fileName string
	funcName string
	lineNum int
	isErrorOrFatal bool
}
//获取调用当前函数的位置


func getLineInfo() (fileName string,funcName string,lineNum int){
	pc,file,line,ok:=runtime.Caller(4)
	if ok {
		funcName=runtime.FuncForPC(pc).Name()
		fileName=file
		lineNum =line
	}
	return
}

/*
 *	使用异步的方式来写日志
 *
 */
func printfLogger(logLevel int,level int,format string,args...interface{}) (chanData *ChanData) {
	//日志等级大于logDebug 不用打印logDeg
	if logLevel > level{
		return
	}
	//打印当前的时间
	timeStr := time.Now().Format("2006-01-02T15:04:05")
	//获取当前的日志的等级信息
	levelInfo:=getLoggerLevelText(level)
	// 获取当前的打印日志的地方
	fileName,funcName,lineNum:=getLineInfo()
	fileName = path.Base(fileName)

	mes:=fmt.Sprintf(format,args...)

	chanData = &ChanData{
		message : mes,
		timeStr :timeStr,
		levelStr :levelInfo,
		fileName : fileName,
		funcName :funcName,
		lineNum :lineNum,
		isErrorOrFatal :false,
	}

	if level == LogLevelError || level ==LogLevelFatal {
		chanData.isErrorOrFatal = true
	}
	return
}



