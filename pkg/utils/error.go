package utils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// CallStack records
func CallStack(depth int) string {
	var callstackStr string
	for i := depth; i >= 0; i-- {
		if pc, file, line, ok := runtime.Caller(i); ok {
			if !strings.Contains(file, "bingo") {
				continue
			}
			fn := runtime.FuncForPC(pc).Name()
			if strings.HasSuffix(fn, "/utils.Recover") {
				break
			}
			fns := strings.Split(fn, ".")
			if len(fns) > 0 {
				fn = fns[len(fns)-1]
			}
			m := file + ":" + strconv.Itoa(line) + " :" + fn
			if callstackStr == "" {
				callstackStr = m
			} else {
				callstackStr += " -> " + m
			}
		}
	}
	return callstackStr
}

// Recover used to catch panic and return an error
func Recover() error {
	if e := recover(); e != nil {
		// callstackStr := CallStack(50)
		// WriteFile("./log/panic", []byte(fmt.Sprintf("%s Panic occurred: %v, call stack: %s\n\n", time.Now(), e, callstackStr)))
		// return log.Fatalf("Panic occurred: %v, call stack: %s\n", e, callstackStr)
		return fmt.Errorf("%s", e)
	}
	return nil
}

const (
	msgNotice  = 0
	msgWarning = 1
	msgError   = 2
)

type Msg struct {
	level int
	err   error
}

func (m *Msg) SetNotice(format string, args ...interface{}) { m.set(msgNotice, format, args...) }
func (m *Msg) SetWarning(format string, args ...interface{}) error {
	return m.set(msgWarning, format, args...)
}
func (m *Msg) SetError(format string, args ...interface{}) error {
	return m.set(msgError, format, args...)
}
func (m *Msg) SetErr(err error) error {
	return m.set(msgError, err.Error())
}
func (m *Msg) set(lvl int, format string, args ...interface{}) error {
	m.level = lvl
	m.err = fmt.Errorf(format, args...)
	return m.err
}

func (m *Msg) Notice() string {
	if m.level == msgNotice && m.err != nil {
		return m.err.Error()
	}
	return ""
}
func (m *Msg) Warning() error {
	if m.level >= msgWarning && m.err != nil {
		return m.err
	}
	return nil
}
func (m *Msg) Error() error {
	if m.level == msgError && m.err != nil {
		return m.err
	}
	return nil
}
