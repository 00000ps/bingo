package format

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	underLine = 4
	twinkle   = 5
	reverse   = 7

	ColorBlack     = 30
	ColorRed       = 31
	ColorGreen     = 32
	ColorYellow    = 33
	ColorBlue      = 34
	ColorMagenta   = 35
	ColorCyan      = 36
	ColorWhite     = 37
	ColorBgBlack   = 40
	ColorBgRed     = 41
	ColorBgGreen   = 42
	ColorBgYellow  = 43
	ColorBgBlue    = 44
	ColorBgMagenta = 45
	ColorBgCyan    = 46
	ColorBgWhite   = 47

	starter = "\x1b["
	ender   = "\x1b[0m"
)

var colorized = true

func SetColorized(c bool) { colorized = c }

func UnderLine(color int, frmt string, a ...interface{}) string {
	return textFmt(underLine, color, frmt, a...)
}
func Reverse(color int, frmt string, a ...interface{}) string {
	return textFmt(reverse, color, frmt, a...)
}
func Twinkle(color int, frmt string, a ...interface{}) string {
	return textFmt(twinkle, color, frmt, a...)
}

func Info(frmt string, a ...interface{}) string    { return textColor(ColorGreen, frmt, a...) }
func Warning(frmt string, a ...interface{}) string { return textColor(ColorMagenta, frmt, a...) }
func Error(frmt string, a ...interface{}) string   { return textColor(ColorRed, frmt, a...) }

func Black(frmt string, a ...interface{}) string   { return textColor(ColorBlack, frmt, a...) }
func Red(frmt string, a ...interface{}) string     { return textColor(ColorRed, frmt, a...) }
func Green(frmt string, a ...interface{}) string   { return textColor(ColorGreen, frmt, a...) }
func Yellow(frmt string, a ...interface{}) string  { return textColor(ColorYellow, frmt, a...) }
func Blue(frmt string, a ...interface{}) string    { return textColor(ColorBlue, frmt, a...) }
func Magenta(frmt string, a ...interface{}) string { return textColor(ColorMagenta, frmt, a...) }
func Cyan(frmt string, a ...interface{}) string    { return textColor(ColorCyan, frmt, a...) }
func White(frmt string, a ...interface{}) string   { return textColor(ColorWhite, frmt, a...) }

func BgBlack(frmt string, a ...interface{}) string   { return textColor(ColorBgBlack, frmt, a...) }
func BgRed(frmt string, a ...interface{}) string     { return textColor(ColorBgRed, frmt, a...) }
func BgGreen(frmt string, a ...interface{}) string   { return textColor(ColorBgGreen, frmt, a...) }
func BgYellow(frmt string, a ...interface{}) string  { return textColor(ColorBgYellow, frmt, a...) }
func BgBlue(frmt string, a ...interface{}) string    { return textColor(ColorBgBlue, frmt, a...) }
func BgMagenta(frmt string, a ...interface{}) string { return textColor(ColorBgMagenta, frmt, a...) }
func BgCyan(frmt string, a ...interface{}) string    { return textColor(ColorBgCyan, frmt, a...) }
func BgWhite(frmt string, a ...interface{}) string   { return textColor(ColorBgWhite, frmt, a...) }

/*func ColorPrintln(s string, i int) { //设置终端字体颜色
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	log.Printf(s)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
}*/
func textColor(color int, strformat string, a ...interface{}) string {
	str := strformat
	if len(a) > 0 {
		str = fmt.Sprintf(strformat, a...)
	}
	if !colorized {
		return str
	}
	if runtime.GOOS == "windows" || color < 30 {
		//	return str
	}
	//return fmt.Sprintf("\033[%dm%s\033[0m", color, str)
	return starter + fmt.Sprintf("%dm", color) + str + ender
}
func textFmt(f, color int, strformat string, a ...interface{}) string {
	str := strformat
	if len(a) > 0 {
		str = fmt.Sprintf(strformat, a...)
	}
	if !colorized {
		return str
	}
	if runtime.GOOS == "windows" {
		//	return str
	}
	if color < 30 {
		return starter + ToString(color, "m") + str + ender
	}
	return starter + ToString(f, ";", color, "m") + str + ender
}

// Clear returns cleared strings without color tags
func Clear(str string) string {
	// return str
	str = strings.Replace(str, "\n", "  ", -1)
	str = strings.Replace(str, "\r", "", -1)
	for i := 0; i < 100; i++ {
		str = strings.Replace(str, ToString(starter, underLine, ";", i, "m"), "", -1)
		str = strings.Replace(str, ToString(starter, twinkle, ";", i, "m"), "", -1)
		str = strings.Replace(str, ToString(starter, reverse, ";", i, "m"), "", -1)
		str = strings.Replace(str, ToString(starter, i, "m"), "", -1)
	}
	return str
}
