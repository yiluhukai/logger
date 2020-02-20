package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileLogger struct {
	level int
	logPath string
	logName string
	file *os.File
	warnFile *os.File
	//存放日志信息的管道
	chanLoggerData  chan *ChanData
	//日志的切分方式
	logSplitType int
	logSplitSize int
	logSplitLastHour int
}



func (logger * FileLogger) init()  {
	seq:=string(filepath.Separator)
	//打开日志文件
	loggerFile:=fmt.Sprintf("%s%s%s.log",logger.logPath,seq,logger.logName)
	file,err:=os.OpenFile(loggerFile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0755)
	if err!=nil{
		panic(fmt.Sprintf("open file %s failed:%v",loggerFile,err))
	}
	logger.file = file
	//打开错误日志文件
	loggerWarnFile:=fmt.Sprintf("%s%s%s.wr.log",logger.logPath,seq,logger.logName)
	file,err=os.OpenFile(loggerWarnFile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0755)
	if err!=nil{
		panic(fmt.Sprintf("open file %s failed:%v",loggerWarnFile,err))
	}
	logger.warnFile = file

	go logger.WriteBackground()

}

func NewFileLogger(config map[string]string) (log Logger, err error) {
	level,ok:=config["level"]
	if !ok{
		err = fmt.Errorf("level is invalid")
		return
	}
	logPath,ok:=config["path_name"]
	if !ok{
		err = fmt.Errorf("logPath doesn't exist")
		return
	}

	logName,ok:=config["file_name"]

	if !ok{
		err = fmt.Errorf("level is invalid")
		return
	}

	chanSize,ok := config["chan_size"]
	if !ok{
		chanSize = "5000"
	}

	logSplitStr,ok := config["log_split_type"]
	var logSize int
	if !ok {
		//默认按时间切割
		logSplitStr = "hour"
	}else{
		if logSplitStr == "size" {
			logSplitSize,ok := config["log_split_size"]
			if !ok {
				//100转成字节
				logSplitSize = "104857600"
			}
			//经代表大小的字符串转成数字
			var errConvert  error
			logSize,errConvert =strconv.Atoi(logSplitSize)
			if errConvert!=nil{
				logSize = 104857600
			}
		}
		//if 等于 "hour"
	}
	//设置日志切分类型
	var logSplitType int
	switch logSplitStr {
	case "hour": logSplitType = LogSplitByHour
	case "size": logSplitType = LogSplitBySize
	default:
		logSplitType = LogSplitByHour
	}

	logger:=&FileLogger{
		level:getLoggerLevel(level),
		logPath:logPath,
		logName:logName,
		logSplitType:logSplitType,
		logSplitSize: logSize,
		logSplitLastHour: time.Now().Hour(),
	}

	//设置管道大小
	size,errCon:= strconv.Atoi(chanSize)
	if errCon!= nil{
		//fmt.Printf("chanSize is invalid:err")
		size = 50000
	}

	logger.chanLoggerData = make(chan *ChanData,size)
	logger.init()
	log =logger
	return
}

func (logger * FileLogger) SetLevel(level int)  {
	if level<LogLevelDebug || level>LogLevelFatal{
		logger.level = LogLevelDebug
		return
	}
	logger.level = level
}

func (logger * FileLogger) LogDebug(format string,args ...interface{})  {

	chanData:=printfLogger(logger.level,LogLevelDebug,format,args...)
	if chanData ==nil{
		return
	}else{
		select {
			case logger.chanLoggerData <- chanData:
			default:
		}
	}
}

func (logger * FileLogger) LogTrace(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelTrace,format,args...)
	if chanData ==nil{
		return
	}else{
		select {
		case logger.chanLoggerData <- chanData:
		default:
		}
	}
}

func (logger * FileLogger) LogInfo(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelInfo,format,args...)
	if chanData ==nil{
		return
	}else{
		select {
		case logger.chanLoggerData <- chanData:
		default:
		}
	}
}

func (logger * FileLogger) LogWarn(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelWarn,format,args...)
	if chanData ==nil{
		return
	}else{
		select {
		case logger.chanLoggerData <- chanData:
		default:
		}
	}
}

func (logger * FileLogger) LogError(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelError,format,args...)
	if chanData ==nil{
		return
	}else{
		select {
		case logger.chanLoggerData <- chanData:
		default:
		}
	}
}

func (logger * FileLogger) LogFatal(format string,args ...interface{})  {
	chanData:=printfLogger(logger.level,LogLevelFatal,format,args...)
	if chanData ==nil{
		return
	}else{
		select {
		case logger.chanLoggerData <- chanData:
		default:
		}
	}
}

func (logger * FileLogger) Close(){
	logger.file.Close()
	logger.warnFile.Close()
}

//检测是否应该切分写日志
func (logger * FileLogger) CheckLogSplit(isErrorOrFatal bool){
	// 获取日志的切分方式
	now:=time.Now()
	curHour:= now.Hour()
	if logger.logSplitType == LogSplitByHour{
		//按时间切分
		if curHour == logger.logSplitLastHour {
			return
		}else{
			logger.LogSplitHour(isErrorOrFatal)
		}
	}else{
		//按大小切分
		logger.LogSplitSize(isErrorOrFatal)
	}
}

func (logger * FileLogger) LogSplitSize(isErrorOrFatal bool){
	//备份文件
	now:=time.Now()
	timeStr:=now.Format("2006-01-02-15-04-05")
	seq := string(filepath.Separator)
	var backupFile, oldFile string
	if isErrorOrFatal{

		fileInfo,err:= logger.warnFile.Stat()
		if err!=nil{
			return
		}
		if fileInfo.Size() < int64(logger.logSplitSize) {
			return
		}
		//关闭当前的日志文件
		logger.warnFile.Close()
		//备份错误日志文件
		backupFile = fmt.Sprintf("%s%s%s-%s.wr.log",logger.logPath,seq,logger.logName,timeStr)
		oldFile =  fmt.Sprintf("%s%s%s.wr.log",logger.logPath,seq,logger.logName)
	}else{
		fileInfo,err:= logger.file.Stat()
		if err!=nil{
			return
		}
		if fileInfo.Size() < int64(logger.logSplitSize) {
			return
		}
		logger.file.Close()
		//备份普通的日志文件
		backupFile = fmt.Sprintf("%s%s%s-%s.log",logger.logPath,seq,logger.logName,timeStr)
		oldFile =  fmt.Sprintf("%s%s%s.log",logger.logPath,seq,logger.logName)
	}

	//重命名文件
	err:=os.Rename(oldFile,backupFile)
	if err !=nil {
		fmt.Printf("备份失败:%v",err)
		//放弃此次备份
		return
	}
	//打开新的文件

	file,err:=os.OpenFile(oldFile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0755)

	if err!=nil{
		//放弃打开
		return
	}
	if isErrorOrFatal{
		logger.warnFile = file

	}else {
		logger.file = file
	}
}

func (logger * FileLogger) LogSplitHour(isErrorOrFatal bool){
	//备份文件
	now:=time.Now()
	timeStr:=now.Format("2006-01-02-15-04")
	seq := string(filepath.Separator)
	var backupFile, oldFile string
	if isErrorOrFatal{
		//关闭当前的日志文件
		logger.warnFile.Close()
		//备份错误日志文件
		backupFile = fmt.Sprintf("%s%s%s-%s.wr.log",logger.logPath,seq,logger.logName,timeStr)
		oldFile =  fmt.Sprintf("%s%s%s.wr.log",logger.logPath,seq,logger.logName)
		//关闭当前文件
	}else{
		logger.file.Close()
		//备份普通的日志文件
		backupFile = fmt.Sprintf("%s%s%s-%s.log",logger.logPath,seq,logger.logName,timeStr)
		oldFile =  fmt.Sprintf("%s%s%s.log",logger.logPath,seq,logger.logName)
	}

	//重命名文件
	err:=os.Rename(oldFile,backupFile)
	if err !=nil {
		//放弃此次备份
		return
	}
	//打开新的文件

	file,err:=os.OpenFile(oldFile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0755)

	if err!=nil{
		//放弃打开
		return
	}
	if isErrorOrFatal{
		logger.warnFile = file

	}else {
		logger.file =file
	}
	logger.logSplitLastHour = now.Hour()
}

//异步写日志的方法
func (logger * FileLogger) WriteBackground(){

	for chanData := range  logger.chanLoggerData{
		file:=logger.file
		if chanData.isErrorOrFatal{
			file = logger.warnFile
		}
		//函数内部会先对文件进行切换再打开
		logger.CheckLogSplit(chanData.isErrorOrFatal)
		log:=fmt.Sprintf("%s %s (%s %s %d) %s",chanData.timeStr,
			chanData.levelStr,chanData.fileName,chanData.funcName,chanData.lineNum,chanData.message)
		fmt.Fprintln(file,log)
	}
}

