package log

func Fatal(format string, a ...interface{}) error   { return nil }
func Error(format string, a ...interface{}) error   { return nil }
func Success(format string, a ...interface{})       {}
func Warning(format string, a ...interface{}) error { return nil }
func Info(format string, a ...interface{})          {}
func Notice(format string, a ...interface{})        {}
func Debug(format string, a ...interface{})         {}
func Trace(format string, a ...interface{})         {}

// func Seperator() { fmt.Println(strings.Repeat("-")) }
