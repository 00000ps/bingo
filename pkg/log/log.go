package log

type (
	// Logger define a interface of log
	Logger interface {
		Fatal(format string, a ...interface{}) error
		Error(format string, a ...interface{}) error
		Success(format string, a ...interface{})
		Warning(format string, a ...interface{}) error
		Info(format string, a ...interface{})
		Notice(format string, a ...interface{})
		Debug(format string, a ...interface{})
		Trace(format string, a ...interface{})
	}

	// level int
)

const (
	// level defines
	traceLevel = iota
	debugLevel
	noticeLevel
	infoLevel
	warningLevel
	successLevel
	errorLevel
	fatalLevel
)
