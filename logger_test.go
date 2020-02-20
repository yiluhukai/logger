package logger

import "testing"

func TestFileLogger(t *testing.T) {
	m:=make(map[string]string,8)
	m["path_name"]=".\\logs"
	m["file_name"]="log"
	m["level"] =getLoggerLevelText(0)
	logger,_:=NewFileLogger(m)

	logger.LogDebug("hello")
	logger.LogInfo("hello%d",10)
	logger.LogError("this is error")
	logger.Close()
}

func TestConsoleLogger(t *testing.T) {
	m:=make(map[string]string,8)
	m["level"] =getLoggerLevelText(0)
	logger,_:=NewConsoleLogger(m)

	logger.LogDebug("hello")
	logger.LogInfo("hello%d",10)
	logger.LogError("this is error")
	logger.Close()
}