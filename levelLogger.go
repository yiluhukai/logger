package logger


//日志的等级
const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

//日志的切分方式


const (
	LogSplitByHour = iota
	LogSplitBySize
)


func getLoggerLevelText(level int) string{
	switch level {
	case LogLevelDebug:
		return "LogLevelDebug"

	case LogLevelTrace:
		return "LogLevelTrace"
	case LogLevelInfo:
		return "LogLevelInfo"
	case LogLevelWarn:
		return "LogLevelWarn"
	case LogLevelError:
		return "LogLevelError"
	case LogLevelFatal:
		return "LogLevelFatal"
	default:
		return "unknown logger"
	}
}


func getLoggerLevel(levelText string) int{
	switch levelText{
	case "LogLevelDebug":
		return LogLevelDebug
	case "LogLevelTrace":
		return LogLevelTrace
	case "LogLevelInfo":
		return LogLevelInfo
	case "LogLevelWarn":
		return LogLevelWarn
	case "LogLevelError":
		return LogLevelError
	case "LogLevelFatal":
		return LogLevelFatal
	default:
		return LogLevelDebug
	}
}



