package format

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 普通占位符
// 占位符     说明                           举例                   输出
// %v      相应值的默认格式。            Printf("%v", people)   {zhangsan}，
// %+v     打印结构体时，会添加字段名     Printf("%+v", people)  {Name:zhangsan}
// %#v     相应值的Go语法表示            Printf("#v", people)   main.Human{Name:"zhangsan"}
// %T      相应值的类型的Go语法表示       Printf("%T", people)   main.Human
// %%      字面上的百分号，并非值的占位符  Printf("%%")            %
// 布尔占位符
// 占位符       说明                举例                     输出
// %t          true 或 false。     Printf("%t", true)       true
// 整数占位符
// 占位符     说明                                  举例                       输出
// %b      二进制表示                             Printf("%b", 5)             101
// %c      相应Unicode码点所表示的字符              Printf("%c", 0x4E2D)        中
// %d      十进制表示                             Printf("%d", 0x12)          18
// %o      八进制表示                             Printf("%d", 10)            12
// %q      单引号围绕的字符字面值，由Go语法安全地转义 Printf("%q", 0x4E2D)        '中'
// %x      十六进制表示，字母形式为小写 a-f         Printf("%x", 13)             d
// %X      十六进制表示，字母形式为大写 A-F         Printf("%x", 13)             D
// %U      Unicode格式：U+1234，等同于 "U+%04X"   Printf("%U", 0x4E2D)         U+4E2D
// 字符串与字节切片
// 占位符     说明                              举例                           输出
// %s      输出字符串表示（string类型或[]byte)   Printf("%s", []byte("Go语言"))  Go语言
// %q      双引号围绕的字符串，由Go语法安全地转义  Printf("%q", "Go语言")         "Go语言"
// %x      十六进制，小写字母，每字节两个字符      Printf("%x", "golang")         676f6c616e67
// %X      十六进制，大写字母，每字节两个字符      Printf("%X", "golang")         676F6C616E67
// 其它标记
// 占位符      说明                             举例          输出
// +      总打印数值的正负号；对于%q（%+q）保证只输出ASCII编码的字符。
//                                            Printf("%+q", "中文")  "\u4e2d\u6587"
// -      在右侧而非左侧填充空格（左对齐该区域）
// #      备用格式：为八进制添加前导 0（%#o）      Printf("%#U", '中')      U+4E2D
//        为十六进制添加前导 0x（%#x）或 0X（%#X），为 %p（%#p）去掉前导 0x；
//        如果可能的话，%q（%#q）会打印原始 （即反引号围绕的）字符串；
//        如果是可打印字符，%U（%#U）会写出该字符的
//        Unicode 编码形式（如字符 x 会被打印成 U+0078 'x'）。
// ' '    (空格)为数值中省略的正负号留出空白（% d）；
//        以十六进制（% x, % X）打印字符串或切片时，在字节之间用空格隔开
// 0      填充前导的0而非空格；对于数字，这会将填充移到正负号之后

// 整数：
// 1 %b    表示为二进制
// 2 %c    该值对应的unicode码值
// 3 %d    表示为十进制
// 4 %o    表示为八进制
// 5 %q    该值对应的单引号括起来的go语法字符字面值，必要时会采用安全的转义表示
// 6 %x    表示为十六进制，使用a-f
// 7 %X    表示为十六进制，使用A-F
// 8 %U    表示为Unicode格式：U+1234，等价于"U+%04X"
// 字符串和[]byte：
// 9  %s    直接输出字符串或者[]byte
//    %q    该值对应的双引号括起来的go语法字符串字面值，必要时会采用安全的转义表示
//    %x    每个字节用两字符十六进制数表示（使用a-f）
//    %X    每个字节用两字符十六进制数表示（使用A-F）
// 指针：
// 10 %p    表示为十六进制，并加上前导的0x
// 布尔值：
// 11 %t    单词true或false
// 浮点数、复数的两个组分：
//    %b    无小数部分、二进制指数的科学计数法，如-123456p-78；参见strconv.FormatFloat
// 12 %e    科学计数法，如-1234.456e+78
// 13 %E    科学计数法，如-1234.456E+78
// 14 %f    有小数部分但无指数部分，如123.456
// 15 %F    等价于%f %g    根据实际情况采用%e或%f格式（以获得更简洁、准确的输出）
// 16 %G    根据实际情况采用%E或%F格式（以获得更简洁、准确的输出）
// 17 %g
// 没有verb %u。整数如果是无符号类型自然输出也是无符号的。类似的，也没有必要指定操作数的尺寸（int8，int64）
// 宽度通过一个紧跟在百分号后面的十进制数指定，如果未指定宽度，则表示值时除必需之外不作填充。精度通过（可能有的）宽度后跟点号后跟的十进制数指定。如果未指定精度，会使用默认精度；如果点号后没有跟数字，表示精度为0。举例如下：
// 1 %f:    默认宽度，默认精度
// 2 %9f    宽度9，默认精度
// 3 %.2f   默认宽度，精度2
// 4 %9.2f  宽度9，精度2
// 5 %9.f   宽度9，精度0
// 宽度和精度格式化控制的是Unicode码值的数量（不同于C的printf，它的这两个因数指的是字节的数量）。两者任一个或两个都可以使用'*'号取代，此时它们的值将被对应的参数（按'*'号和verb出现的顺序，即控制其值的参数会出现在要表示的值前面）控制，这个操作数必须是int类型。
// 对于大多数类型的值，宽度是输出的最小字符数，如果必要是会用空格填充。对于字符串，宽度是输出字符数目的最低数量，如果必要会截断字符串。
// 对于整数，宽度和精度都设置输出总长度。采用精度时表示右对齐并用0填充，而宽度默认表示用空格填充。
// 对于浮点数，宽度设置输出总长度；精度设置小数部分长度（如果有的话），除了%g/%G，此时精度设置总的数字个数。例如，对数字123.45，格式%6.2f 输出123.45；格式%.4g输出123.5。%e和%f的默认精度是6，%g的默认精度是可以将该值区分出来需要的最小数字个数。
// 对复数，宽度和精度会分别用于实部和虚部，结果用小括号包裹。因此%f用于1.2+3.4i输出(1.200000+3.400000i)。
// 其它flag：
// 1 +    总是输出数值的正负号；对%q（%+q）会生成全部是ASCII字符的输出（通过转义）；
// 2 -    在输出右边填充空白而不是默认的左边（即从默认的右对齐切换为左对齐）；
// 3 #    切换格式：
// 4      八进制数前加0（%#o），十六进制数前加0x（%#x）或0X（%#X），指针去掉前面的0x（%#p）；
// 5      对%q（%#q），如果strconv.CanBackquote返回真会输出反引号括起来的未转义字符串；
// 6      对%U（%#U），如果字符是可打印的，会在输出Unicode格式、空格、单引号括起来的go字面值；
// 7 ' '  对数值，正数前加空格而负数前加负号；
// 8      对字符串采用%x或%X时（% x或% X）会给各打印的字节之间加空格；
// 9 0    使用0而不是空格填充，对于数值类型会把填充的0放在正负号后面；
// 显式指定参数索引：
// 在Printf、Sprintf、Fprintf三个函数中，默认的行为是对每一个格式化verb依次对应调用时成功传递进来的参数。但是，紧跟在verb之前的[n]符号表示应格式化第n个参数（索引从1开始）。同样的在'*'之前的[n]符号表示采用第n个参数的值作为宽度或精度。在处理完方括号表达式[n]后，除非另有指示，会接着处理参数n+1，n+2……（就是说移动了当前处理位置）。例如：
// 1 fmt.Sprintf("%[2]d %[1]d\n", 11, 22)
// 会生成"22 11"，而：
// 1 fmt.Sprintf("%[3]*.[2]*[1]f", 12.0, 2, 6),
// 等价于：
// 1 fmt.Sprintf("%6.2f", 12.0),
// 会生成" 12.00"。因为显式的索引会影响随后的verb，这种符号可以通过重设索引用于多次打印同一个值：
// 1 fmt.Sprintf("%d %d %#[1]x %#x", 16, 17)
// 会生成"16 17 0x10 0x11"

// 格式化错误：
// 如果给某个verb提供了非法的参数，如给%d提供了一个字符串，生成的字符串会包含该问题的描述，如下所例：
//  1 错误的类型或未知的verb：%!verb(type=value)
//  2     Printf("%d", hi):          %!d(string=hi)
//  3 太多参数（采用索引时会失效）：%!(EXTRA type=value)
//  4     Printf("hi", "guys"):      hi%!(EXTRA string=guys)
//  5 太少参数: %!verb(MISSING)
//  6     Printf("hi%d"):            hi %!d(MISSING)
//  7 宽度/精度不是整数值：%!(BADWIDTH) or %!(BADPREC)
//  8     Printf("%*s", 4.5, "hi"):  %!(BADWIDTH)hi
//  9     Printf("%.*s", 4.5, "hi"): %!(BADPREC)hi
// 10 没有索引指向的参数：%!(BADINDEX)
// 11     Printf("%*[2]d", 7):       %!d(BADINDEX)
// 12     Printf("%.[2]d", 7):       %!d(BADINDEX)
// 所有的错误都以字符串"%!"开始，有时会后跟单个字符（verb标识符），并以加小括弧的描述结束。
// 如果被print系列函数调用时，Error或String方法触发了panic，fmt包会根据panic重建错误信息，用一个字符串说明该panic经过了fmt包。例如，一个String方法调用了panic("bad")，生成的格式化信息差不多是这样的：
// %!s(PANIC=bad)
// %!s指示表示错误（panic）出现时的使用的verb。

const (
	valErr            = -100
	valNotSet         = -1
	valDefaultWidth   = -1
	valDefaultMinSpec = -1
	valDefaultSpec    = 6
	valDefaultBase    = 10
)

var (
	errNil       = errors.New("%!MISSING")
	errVerb      = errors.New("%!MISSING VERB")
	errIndex     = errors.New("%!BADINDEX") // %!f(BADINDEX)
	errWidth     = errors.New("%!BAD WIDTH")
	errSpec      = errors.New("%!BAD SPEC")
	errWidthSpec = errors.New("%!BAD WIDTH & SPEC")
	errPara      = errors.New("%!BAD PARAMETER")
	errParaMiss  = errors.New("%!MISSING PARAMETER")

	tds []ts
)

type ts struct {
	fn   string
	s, e time.Time
}
type flag struct {
	prefix []byte
	suffix []byte
	verb   byte

	indexStr string
	index    int // '[n]'  eg
	width    int //
	prec     int //        eg  %f:默认宽度，默认精度;  %9f:宽度9，默认精度;  %.2f:默认宽度，精度2;  %9.2f:宽度9，精度2;  %9.f:宽度9，精度0
	base     int

	upper    bool
	ws       bool
	dot      bool
	indexing bool
	widthMod bool
	precMod  bool
	zero     bool // '0'    eg  使用0而不是空格填充，对于数值类型会把填充的0放在正负号后面
	space    bool // ' '    eg  对数值，正数前加空格而负数前加负号;  对字符串采用%x或%X时（% x或% X）会给各打印的字节之间加空格
	minus    bool // '-'    eg  在输出右边填充空白而不是默认的左边（即从默认的右对齐切换为左对齐）
	plus     bool // '+'    eg  总是输出数值的正负号；对%q（%+q）会生成全部是ASCII字符的输出（通过转义）；
	sharp    bool // '#'    eg  切换格式.  八进制数前加0（%#o），十六进制数前加0x（%#x）或0X（%#X），指针去掉前面的0x（%#p）;  对%q（%#q），如果strconv.CanBackquote返回真会输出反引号括起来的未转义字符串;  对%U（%#U），如果字符是可打印的，会在输出Unicode格式、空格、单引号括起来的go字面值
}

// // flags placed in a separate struct for easy clearing.
// type fmtFlags struct {
// 	// intbuf is large enough to store %b of an int64 with a sign and
// 	// avoids padding at the end of the struct on 32 bit architectures.
// 	intbuf [68]byte
// 	buf    *buffer

// 	wid  int // width
// 	prec int // precision

// 	widPresent  bool
// 	precPresent bool
// 	minus       bool
// 	plus        bool
// 	sharp       bool
// 	space       bool
// 	zero        bool

// 	// For the formats %+v %#v, we set the plusV/sharpV flags
// 	// and clear the plus/sharp flags since %+v and %#v are in effect
// 	// different, flagless formats set at the top level.
// 	plusV  bool
// 	sharpV bool
// }

// // Println is
// func Println(f string, a ...interface{}) { fmt.Println(Sprintf(f, a...)) }

// // Printf is
// func Printf(f string, a ...interface{}) { Sprintf(f, a...) }

// // func Printf(f string, a ...interface{}) { fmt.Printf(Sprintf(f, a...)) }

// // Sprintf is a memory safe string combiner
// func Sprintf(f string, a ...interface{}) string {
// 	st := time.Now()
// 	control := ""
// 	control = fmt.Sprintf(f, a...)
// 	cCost := time.Since(st)
// 	// debug("\n\n\n")
// 	dbg := &dg{on: true}
// 	// dbg.debug("\n\n\n[Sprintf] %s", control)

// 	st = time.Now()
// 	result := dbg.sprintf(f, a...)
// 	tCost := time.Since(st)
// 	inc := (float32(tCost) - float32(cCost)) * 100 / float32(cCost)

// 	if false {
// 		fmt.Printf("\n")
// 		for i, ts := range tds {
// 			c := ts.e.Sub(ts.s)
// 			d := (float32(c)) * 100 / float32(cCost)
// 			fmt.Printf("-- %2d: %6.1f%% %10s %-s\n", i, d, c, ts.fn)
// 		}
// 	}
// 	tds = []ts{}

// 	if control == "" || result == control {
// 		dbg.msg = []string{}
// 		head := "[Sprintf] \x1b[32mPASS\x1b[0m \x1b[31m%+10.0f%%\x1b[0m %10s=%10s-%8s fool:\x1b[33m%5t\x1b[0m len:\x1b[34m%-2d\x1b[0m RESULT:\x1b[32m%q\x1b[0m format:\x1b[34m%s\x1b[0m args:\x1b[34m%#v\x1b[0m\n"
// 		dbg.debug(head, inc, tCost-cCost, tCost, cCost, dbg.fool, len(a), result, f, a)

// 		// fmt.Printf(head, inc, tCost-cCost, tCost,cCost, dbg.fool, len(a), result, f, a)
// 	} else {
// 		head := "[Sprintf] \x1b[31mFAIL\x1b[0m \x1b[31m%+10.0f%%\x1b[0m %10s=%10s-%8s fool:\x1b[33m%5t\x1b[0m len:\x1b[34m%-2d\x1b[0m RESULT:\x1b[32m%q\x1b[0m format:\x1b[34m%s\x1b[0m  args:\x1b[34m%#v\x1b[0m\n\n\n"
// 		dbg.debug("[Sprintf] format: \x1b[34m%s\x1b[0m; length: \x1b[34m%d\x1b[0m;  args: \x1b[34m%#v\x1b[0m", f, len(a), a)
// 		dbg.debug("[Sprintf] CONTRL: \x1b[32m%q\x1b[0m", control)
// 		dbg.debug("[Sprintf] RESULT: \x1b[31m%q\x1b[0m", result)
// 		dbg.debug(head, inc, tCost-cCost, tCost, cCost, dbg.fool, len(a), result, f, a)

// 		// fmt.Printf("[Sprintf] format: \x1b[34m%s\x1b[0m; length: \x1b[34m%d\x1b[0m;  args: \x1b[34m%#v\x1b[0m\n", f, len(a), a)
// 		// fmt.Printf("[Sprintf] CONTRL: \x1b[32m%q\x1b[0m\n", control)
// 		// fmt.Printf("[Sprintf] RESULT: \x1b[31m%q\x1b[0m\n", result)
// 		// fmt.Printf(head, inc, tCost-cCost, tCost, cCost, dbg.fool, len(a), result, f, a)
// 	}
// 	dbg.record()
// 	// os.Exit(1)
// 	return result
// }
func (dbg *dg) sprintf(f string, a ...interface{}) string {
	stt := time.Now()
	// dbg.debug("[sprintf] format:\x1b[34m" + f + "\x1b[0m")
	// dbg.debug("[sprintf] length:\x1b[34m%d\x1b[0m args:\x1b[34m%#v\x1b[0m", len(a), a)

	st1 := time.Now()
	perFlag := "*BINGO<percent]"
	f = strings.Replace(f, "%%", perFlag, -1)
	tds = append(tds, ts{fn: "sprintf-strings.Replace1", s: st1, e: time.Now()})
	st2 := time.Now()
	slice := strings.Split(f, "%")
	tds = append(tds, ts{fn: "sprintf-strings.Split", s: st2, e: time.Now()})
	st3 := time.Now()
	firstIgnore := !strings.HasPrefix(f, "%") || slice[0] == ""
	tds = append(tds, ts{fn: "sprintf-strings.HasPrefix", s: st3, e: time.Now()})
	// dbg.debug("[sprintf] length:\x1b[34m%d\x1b[0m slice:\x1b[34m%#v\x1b[0m, ignoreFirst:\x1b[34m%t\x1b[0m", len(slice), slice, firstIgnore)
	// missArgs := len(strSlice) == len(a)
	st4 := time.Now()
	b := bytes.NewBuffer(make([]byte, 0, 1000))
	tds = append(tds, ts{fn: "sprintf-bytes.NewBuffer", s: st4, e: time.Now()})
	for i, s := range slice {
		// dbg.debug(strings.Repeat("-", 100))
		if s == "" {
			continue
		}
		idx := i
		if firstIgnore {
			idx = i - 1
		}
		// dbg.debug("[sprintf] \x1b[34m%d\x1b[0m >>>> \x1b[33m%q\x1b[0m", i, s)
		if i == 0 && firstIgnore {
			st41 := time.Now()
			b.WriteString(s)
			tds = append(tds, ts{fn: "sprintf-WriteString", s: st41, e: time.Now()})
		} else if v, e := dbg.print(string(s), idx, a); e == nil {
			// dbg.debug("[sprintf] \x1b[34m%d\x1b[0m <<<< \x1b[33m\"%%%s\"\x1b[0m -----> \x1b[33m%s\x1b[0m", i, s, v)
			st42 := time.Now()
			b.Write(v)
			tds = append(tds, ts{fn: "sprintf-Write", s: st42, e: time.Now()})
		} else {
			// dbg.debug("[sprintf] error: \x1b[31m%s\x1b[0m", e)
		}
	}

	// dbg.debug(strings.Repeat("-", 100))
	// return b.String()
	st5 := time.Now()
	s := strings.Replace(b.String(), perFlag, "%%", -1)
	tds = append(tds, ts{fn: "sprintf-strings.Replace2", s: st5, e: time.Now()})
	tds = append(tds, ts{fn: "sprintf", s: stt, e: time.Now()})
	return s
}
func (dbg *dg) print(str string, idx int, args []interface{}) (s []byte, e error) {
	stt := time.Now()
	if len(str) == 0 {
		e = errNil
		return
	}

	var (
		prefix  []byte
		cutHead = true
		a       = []byte(str)
		head    = a[0]
		f       = flag{
			index: valNotSet,
			base:  valDefaultBase,
			width: valDefaultWidth,
			prec:  valDefaultSpec,
		}
	)

	// debug("\t[print] %d str:%q", idx, str)
	// 提取转义字符
	// if len(a) == 1{
	// 	f.verb = a[0]
	// } else {

	for i, r := range a {
		switch r {
		case 'b', 'c', 'd', 'e', 'E', 'f', 'F', 'g', 'G', 'o', 'p', 'q', 's', 't', 'T', 'U', 'v', 'x', 'X':
			f.verb = r
			prefix = a[:i]
			if i+1 <= len(a) {
				f.suffix = a[i+1:]
			}
			goto tagA
		}
	}

	// for i:=0;i<len(a);i++ {
	// 	switch a[i] {
	// 	case 'b', 'c', 'd', 'e', 'E', 'f', 'F', 'g', 'G', 'o', 'p', 'q', 's', 't', 'T', 'U', 'v', 'x', 'X':
	// 		f.verb = a[i]
	// 		prefix = a[:i]
	// 		if i+1 <= len(a) {
	// 			f.suffix = a[i+1:]
	// 		}
	// 		goto tagA
	// 	}
	// }
	// }
tagA:
	// dbg.debug("\t[print] \x1b[34m%d\x1b[0m prefix:\x1b[34m%q\x1b[0m, suffix:\x1b[34m%q\x1b[0m, head:\x1b[34m%c\x1b[0m, verb:\x1b[34m%c\x1b[0m, f:\x1b[34m%+v\x1b[0m", idx, prefix, f.suffix, head, f.verb, f)

	switch f.verb {
	// 整数
	case 'b':
		f.base = 2
	case 'o':
		f.base = 8
	case 'd':
		f.base = 10
	case 'X':
		f.base = 16
		f.upper = true
	case 'x':
		f.base = 16
	case 'U': // fmt.Printf("%U\n", 0x4E2D) // U+4E2D
		f.base = 16
		f.upper = true
	case 'q': // fmt.Printf("%q\n", 0x4E2D) // '中'
		f.base = 16
	case 'c': // fmt.Printf("%c\n", 0x4E2D) //  中

	// case 'F':
	// 	slice = strconv.AppendFloat(slice, vs.Float(), 'f', f.prec, 64)
	// case 'v':
	// 	slice = strconv.AppendFloat(slice, vs.Float(), 'f', valDefaultWidth, 64)

	// 浮点
	case 'E':
		f.upper = true
	case 'e':
	case 'F':
		f.verb = 'f'
		f.upper = true
	case 'f':
	case 'g', 'G', 'v':
		f.prec = valDefaultMinSpec

	// string / []byte
	case 's':

	// bool
	case 't':

	case 'T':

	default:
		e = errVerb
		return
	}

	tds = append(tds, ts{fn: "print-initttt", s: stt, e: time.Now()})
	// 分析格式化字符串
	if len(prefix) > 0 {
		st1 := time.Now()
		f.indexing = bytes.HasSuffix(prefix, []byte("]"))
		tds = append(tds, ts{fn: "print-bytes.HasSuffix", s: st1, e: time.Now()})
		st2 := time.Now()
		f.ws = bytes.Contains(prefix, []byte("."))
		tds = append(tds, ts{fn: "print-bytes.Contains", s: st2, e: time.Now()})
		f.dot = f.ws

		// 提取起始字符
		switch head {
		case '+':
			f.plus = true
		case '-':
			f.minus = true
			f.zero = false
		case ' ':
			f.space = true
		case '0':
			f.zero = !f.minus
		case '#':
			f.sharp = true

		case '[':
			f.indexing = true
			cutHead = false
		case '.', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			f.ws = true
			cutHead = false

		case 'b', 'c', 'd', 'e', 'E', 'f', 'F', 'g', 'G', 'o', 'p', 'q', 's', 't', 'T', 'U', 'v', 'x', 'X':
		default:
			e = errPara
			return
		}

		if cutHead {
			prefix = prefix[1:]
		}

		// 分析格式化字符串
		if len(prefix) > 0 {
			// 分析格式化内容
			if f.ws || f.indexing {
				// if bytes.Contains(msg, []byte(".")) {
				// * 分析是否指定了width和spec
				f.width, f.prec, f.index, e = dbg.parseWidthSpec(f, prefix, args)
				if e != nil {
					return
				}
			} else if len(prefix) == 0 {
			} else if v, e := strconv.Atoi(string(prefix)); e == nil {
				f.width = v
			} else {
				e = errPara
			}
		}
	}

	st2 := time.Now()
	dp := dbg.printf(f, idx, args)
	tds = append(tds, ts{fn: "print-printf", s: st2, e: time.Now()})
	st3 := time.Now()
	s = append(dp, f.suffix...)
	tds = append(tds, ts{fn: "print-append", s: st3, e: time.Now()})
	tds = append(tds, ts{fn: "print", s: stt, e: time.Now()})
	return
}
func (dbg *dg) parseWidthSpec(f flag, m []byte, args []interface{}) (width, prec, index int, e error) {
	stt := time.Now()
	width = f.width
	prec = f.prec
	index = f.index

	st1 := time.Now()
	a := bytes.Split(m, []byte("."))
	tds = append(tds, ts{fn: "parseWidthSpec-strings.Split", s: st1, e: time.Now()})
	st2 := time.Now()
	// dbg.debug("\t[parseWidthSpec] format:\x1b[34m%q\x1b[0m, length:\x1b[34m%d\x1b[0m, slice:\x1b[34m%#v\x1b[0m", m, len(a), a)
	tds = append(tds, ts{fn: "parseWidthSpec-dbg.debug", s: st2, e: time.Now()})
	switch len(a) {
	case 2:
		for i, asl := range a {
			v := 0

			// %.2f %9.f
			if len(asl) == 0 {
				// dbg.debug("\t[parseWidthSpec] \x1b[34m%d\x1b[0m %%.2f %%9.f", i)
				if i == 1 {
					v = 0
				} else if i == 0 {
					v = valDefaultWidth
				}
			} else if v, e = strconv.Atoi(string(asl)); e == nil {
				// %9.2f
				// dbg.debug("\t[parseWidthSpec] \x1b[34m%d\x1b[0m %%9.2f", i)
				// e = errWidthSpec
				// return
			} else if bytes.HasPrefix(asl, []byte("[")) && bytes.Contains(asl, []byte("]*")) {
				// %[2]*.[1]*[3]f -> [2]*   [1]*[2]
				// %[2]*.[1]*f    -> [2]*   [1]*
				// dbg.debug("\t[parseWidthSpec] \x1b[34m%d\x1b[0m %%[2]*.[1]*[3]f", i)
				startIdx := bytes.Index(asl, []byte("*"))
				if startIdx < 3 {
					e = errIndex
					return
				}
				var ti int
				if ti, e = dbg.parseIndex(asl[:startIdx], args); e == nil {
					tar := args[ti]
					switch tar.(type) {
					case int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
						log.Printf("args[%d]=%#v, type:%#v\n", i-1, tar, tar.(interface{}))
						v = tar.(int)
					case int:
						v = tar.(int)
					}
				} else {
					return
				}
				if bytes.HasSuffix(asl, []byte("]")) {
					// [1]*[2]   index:$2
					if index, e = dbg.parseIndex(asl[startIdx+1:], args); e != nil {
						return
					}
				}
			} else if i == 1 && bytes.HasSuffix(asl, []byte("]")) && bytes.Contains(asl, []byte("[")) {
				il := bytes.Index(asl, []byte("["))
				ir := bytes.Index(asl, []byte("]"))
				if v, e = strconv.Atoi(string(asl[:il])); e != nil {
					return
				}

				if index, e = dbg.parseIndex(asl[il+1:ir], args); e != nil {
					return
				}
				// dbg.debug("\t[parseWidthSpec] \x1b[34m%d\x1b[0m %%6.2[2]f  --> %q|%q ==> %d|%d", i,f[:il],f[il+1:ir], v, index)
			} else {
				// dbg.debug("\t[parseWidthSpec] \x1b[34m%d\x1b[0m error 2", i)
				e = errWidthSpec
				return
			}

			// dbg.debug("\t[parseWidthSpec] \x1b[34m%d\x1b[0m setv=%d", i, v)
			if i == 0 {
				width = v
			} else if i == 1 {
				prec = v
			}
		}
	case 1:
		if bytes.Contains(m, []byte("]")) {
			// * 只有索引 %[2]d
			// dbg.debug("\t[parseWidthSpec] %%[2]d")
			if index, e = dbg.parseIndex(m, args); e != nil {
				return
			}
		} else if width, e = strconv.Atoi(string(m)); width <= 0 || e != nil {
			// * 只有宽度 %9d
			// dbg.debug("\t[parseWidthSpec] %%9d")
			e = errWidth
			return
		}
	default:
		// dbg.debug("\t[parseWidthSpec] error 1")
		e = errWidthSpec
		return
	}

	// dbg.debug("\t[parseWidthSpec] width:\x1b[34m%d\x1b[0m, prec:\x1b[34m%d\x1b[0m, index:\x1b[34m%d\x1b[0m", width, prec, index)

	tds = append(tds, ts{fn: "parseWidthSpec", s: stt, e: time.Now()})
	return

}
func (dbg *dg) parseIndex(idx []byte, args []interface{}) (int, error) {
	stt := time.Now()
	if idx[0] == byte('[') {
		idx = idx[1:]
	}
	l := len(idx)
	if idx[l-1] == byte(']') {
		idx = idx[:l-1]
	}
	// idx = strings.TrimPrefix(idx, "[")
	// idx = strings.TrimSuffix(idx, "]")

	i, e := strconv.Atoi(string(idx))
	// dbg.debug("\t[parseIndex] %s->%d len:%d", idx, i, len(args))
	if e == nil && i > 0 && len(args) >= i {
		tds = append(tds, ts{fn: "parseIndex", s: stt, e: time.Now()})
		return i - 1, nil
	}

	tds = append(tds, ts{fn: "parseIndex", s: stt, e: time.Now()})
	return valErr, errIndex
}
func (dbg *dg) printf(f flag, idx int, args []interface{}) []byte {
	stt := time.Now()
	// dbg.debug("\t[printf] \x1b[34m%d\x1b[0m f.index:\x1b[33m%d\x1b[0m len:%d f:\x1b[34m%+v\x1b[0m", idx, f.index, len(args), f)
	if f.index < 0 {
		f.index = idx
	}
	// dbg.debug("\t[printf] \x1b[34m%d\x1b[0m f.index:\x1b[33m%d\x1b[0m args:\x1b[32m%#v\x1b[0m, f:\x1b[34m%+v\x1b[0m", idx, f.index, args[f.index], f)
	o := args[f.index]
	b := dbg.toBytes(f, o)
	var (
		prefix []byte
		suffix []byte
		// buf    = bytes.NewBuffer(make([]byte, len(b)+100))
	)

	// plus         bool // '+'    eg  总是输出数值的正负号；对%q（%+q）会生成全部是ASCII字符的输出（通过转义）；
	// sharp        bool // '#'    eg  切换格式.  八进制数前加0（%#o），十六进制数前加0x（%#x）或0X（%#X），指针去掉前面的0x（%#p）;  对%q（%#q），如果strconv.CanBackquote返回真会输出反引号括起来的未转义字符串;  对%U（%#U），如果字符是可打印的，会在输出Unicode格式、空格、单引号括起来的go字面值
	// zero     bool // '0'    eg  使用0而不是空格填充，对于数值类型会把填充的0放在正负号后面
	// space    bool // ' '    eg  对数值，正数前加空格而负数前加负号;  对字符串采用%x或%X时（% x或% X）会给各打印的字节之间加空格
	// minus bool // '-'    eg  在输出右边填充空白而不是默认的左边（即从默认的右对齐切换为左对齐）

	if f.verb == 'U' {
		prefix = []byte("U+")
	}

	if f.plus {
		switch f.verb {
		case 'c', 'U':
		case 'b', 'd', 'o', 'x', 'X', 'e', 'E', 'f', 'F', 'g', 'G':
			switch o.(type) {
			case uint, uint8, uint16, uint32, uint64:
				// buf.WriteRune('+')
				prefix = []byte("+")
			case int, int8, int16, int32, int64:
				if reflect.ValueOf(o).Int() > 0 {
					// buf.WriteRune('+')
					prefix = []byte("+")
				} else {
					prefix = []byte("-")
					b = b[1:]
				}
			case float32, float64:
				if reflect.ValueOf(o).Float() > 0 {
					// buf.WriteRune('+')
					prefix = []byte("+")
				} else {
					prefix = []byte("-")
					b = b[1:]
				}
				// case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				// case 'b', 'c', 'd', 'e', 'E', 'f', 'F', 'g', 'G', 'o', 'p', 'q', 's', 't', 'T', 'U', 'v', 'x', 'X':
				// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			}
		case 'q':
			prefix = []byte("'\\u")
			suffix = []byte("'")
		}
	}
	if f.sharp {
		switch f.verb {
		case 'o':
			// buf.WriteRune('0')
			prefix = []byte("0")
		case 'x':
			// buf.WriteString("0x")
			prefix = []byte("0x")
		case 'X':
			// buf.WriteString("0X")
			prefix = []byte("0X")
		case 'p':
			b = bytes.TrimLeft(b, "0x")
		}
	}

	paddlen := 0
	if f.space {
		switch f.verb {
		case 'c', 'q', 'U':
		default:
			paddlen = 1
		}
	}
	if paddlen == 1 || (f.width != valDefaultWidth && f.width > 0) {
		currentlen := len(string(b)) + len(string(prefix))

		// paddlen = 0
		if currentlen < f.width {
			paddlen = f.width - currentlen
		}

		// if currentlen < f.width {
		if f.minus {
			prefix = append(prefix, b...)
			b = bytes.Repeat([]byte(" "), paddlen)
		} else {
			if f.zero {
				b = append(bytes.Repeat([]byte("0"), paddlen), b...)
			} else {
				b = append(prefix, b...)
				prefix = bytes.Repeat([]byte(" "), paddlen)
			}
		}
	}

	if f.upper {
		b = bytes.ToUpper(b)
	}

	// buf.Write(b)
	// buf.Write(suffix)
	// return buf.Bytes()
	r := append(f.prefix, prefix...)
	r = append(r, b...)
	r = append(r, suffix...)
	tds = append(tds, ts{fn: "printf", s: stt, e: time.Now()})
	return r
}
func (dbg *dg) toBytes(f flag, a interface{}) []byte {
	stt := time.Now()
	slice := []byte{}

	var (
		vint                int64
		vuint               uint64
		vfloat              float64
		bint, buint, bfloat bool
	)

	switch a.(type) {
	case bool:
		if a.(bool) {
			slice = []byte("true")
		} else {
			slice = []byte("false")
		}

	case int:
		vint = int64(a.(int))
		bint = true
	case int8:
		vint = int64(a.(int8))
		bint = true
	case int16:
		vint = int64(a.(int16))
		bint = true
	case int32:
		vint = int64(a.(int32))
		bint = true
	case int64:
		vint = a.(int64)
		bint = true
	case uint:
		vuint = uint64(a.(uint))
		buint = true
	case uint8:
		vuint = uint64(a.(uint8))
		buint = true
	case uint16:
		vuint = uint64(a.(uint16))
		buint = true
	case uint32:
		vuint = uint64(a.(uint32))
		buint = true
	case uint64:
		vuint = a.(uint64)
		buint = true

	case float32:
		vfloat = float64(a.(float32))
		bfloat = true
	case float64:
		vfloat = a.(float64)
		bfloat = true

	case string:
		v := a.(string)
		// v := vs.String()
		switch f.verb {
		case 's', 'v':
			slice = []byte(v)
		case 'q':
			slice = strconv.AppendQuote(slice, v)
		case 'x', 'X':
			tb := []byte(v)
			for _, b := range tb {
				slice = strconv.AppendInt(slice, int64(b), 16)
			}
		}
	case []byte:
		v := a.([]byte)
		switch f.verb {
		case 'x', 'X':
			for _, b := range v {
				slice = strconv.AppendInt(slice, int64(b), 16)
			}
		default:
			slice = v
		}
	default:
		goto fool
	}

	if bint {
		switch f.verb {
		case 's':
			vs := reflect.ValueOf(a)
			method := vs.MethodByName("String")
			if method.IsValid() {
				ret := method.Call(make([]reflect.Value, 0))
				slice = []byte(ret[0].String())
			}
		case 'o', 'b', 'd', 'x', 'X', 'v':
			slice = strconv.AppendInt(slice, vint, f.base)
		case 'U':
			slice = strconv.AppendInt(slice, vint, 16)
			if f.sharp {
				slice = append(slice, byte(' '))
				slice = strconv.AppendQuoteRune(slice, rune(vint))
			}
		case 'q':
			if f.plus {
				vs := reflect.ValueOf(a)
				slice = strconv.AppendInt(slice, vs.Int(), f.base)
			} else {
				slice = strconv.AppendQuoteRune(slice, rune(vint))
			}
		case 'c':
			slice = bytes.Trim(strconv.AppendQuoteRune(slice, rune(vint)), "'")
		default:
		}
	} else if buint {
		switch f.verb {
		case 's':
			vs := reflect.ValueOf(a)
			method := vs.MethodByName("String")
			if method.IsValid() {
				ret := method.Call(make([]reflect.Value, 0))
				slice = []byte(ret[0].String())
			}
		case 'o', 'b', 'd', 'x', 'X', 'v':
			slice = strconv.AppendUint(slice, vuint, f.base)
		case 'U':
			slice = strconv.AppendUint(slice, vuint, 16)
			if f.sharp {
				slice = append(slice, byte(' '))
				slice = strconv.AppendQuoteRune(slice, rune(vuint))
			}
		case 'q':
			if f.plus {
				slice = strconv.AppendUint(slice, vuint, 16)
			} else {
				slice = strconv.AppendQuoteRune(slice, rune(vuint))
			}
		case 'c':
			slice = bytes.Trim(strconv.AppendQuoteRune(slice, rune(vuint)), "'")
		default:
		}
	} else if bfloat {
		switch f.verb {
		case 'v':
			if f.dot {
				if f.prec == 0 {
					slice = strconv.AppendFloat(slice, vfloat, 'e', 0, 64)
				} else {
					slice = strconv.AppendFloat(slice, vfloat, 'f', 0, 64)
				}
			} else {
				slice = strconv.AppendFloat(slice, vfloat, 'f', f.prec, 64)
			}
		default:
			slice = strconv.AppendFloat(slice, vfloat, f.verb, f.prec, 64)
		}
	}

	tds = append(tds, ts{fn: "toBytes*", s: stt, e: time.Now()})
	return slice
fool:
	// dbg.debug("\t[toBytes] <<<< \x1b[32m%s\x1b[0m by fool", slice)
	dbg.fool = true
	slice = append(slice, []byte(fmt.Sprintf("%"+string(f.verb), a))...)
	tds = append(tds, ts{fn: "toBytes", s: stt, e: time.Now()})
	return slice

	// 	slice = []byte{}

	// 	stt = time.Now()
	// 	vs := reflect.ValueOf(a)
	// 	// dbg.debug("\t[toBytes] type:\x1b[32m%s\x1b[0m, args:\x1b[32m%#v\x1b[0m, verb:\x1b[34m%c\x1b[0m, f:\x1b[34m%+v\x1b[0m", vs.Kind(), a, f.verb, f)

	// 	switch vs.Kind() {
	// 	case reflect.Bool:
	// 		// debug("------> 0 bool")
	// 		slice = strconv.AppendBool(slice, a.(bool))
	// 		// b.Write()
	// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 		// debug("------> 1 int")
	// 		switch f.verb {
	// 		case 's':
	// 			method := vs.MethodByName("String")
	// 			if method.IsValid() {
	// 				ret := method.Call(make([]reflect.Value, 0))
	// 				slice = []byte(ret[0].String())
	// 			}
	// 		case 'o':
	// 			// slice = []byte("0")
	// 			fallthrough
	// 		case 'b', 'd', 'x', 'X', 'v':
	// 			slice = strconv.AppendInt(slice, vs.Int(), f.base)
	// 		case 'U':
	// 			slice = strconv.AppendInt(slice, vs.Int(), 16)
	// 			if f.sharp {
	// 				slice = append(slice, byte(' '))
	// 				slice = strconv.AppendQuoteRune(slice, rune(a.(int)))
	// 			}
	// 		case 'q':
	// 			if f.plus {
	// 				slice = strconv.AppendInt(slice, vs.Int(), f.base)
	// 			} else {
	// 				slice = strconv.AppendQuoteRune(slice, rune(a.(int)))
	// 			}
	// 		case 'c':
	// 			slice = bytes.Trim(strconv.AppendQuoteRune(slice, rune(a.(int))), "'")
	// 		default:
	// 			goto fool
	// 		}
	// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// 		// debug("------> 2 uint")
	// 		// slice = strconv.AppendUint(slice, vs.Uint(), f.base)
	// 		switch f.verb {
	// 		case 's':
	// 			method := vs.MethodByName("String")
	// 			if method.IsValid() {
	// 				ret := method.Call(make([]reflect.Value, 0))
	// 				slice = []byte(ret[0].String())
	// 			}
	// 		case 'o':
	// 			slice = []byte("0")
	// 			fallthrough
	// 		case 'b', 'd', 'x', 'X', 'v':
	// 			slice = strconv.AppendUint(slice, vs.Uint(), f.base)
	// 		case 'U':
	// 			slice = strconv.AppendUint([]byte("U+"), vs.Uint(), 16)
	// 		case 'q':
	// 			if f.plus {
	// 				slice = strconv.AppendUint(slice, vs.Uint(), 16)
	// 			} else {
	// 				slice = strconv.AppendQuoteRune(slice, rune(vs.Uint()))
	// 			}
	// 		case 'c':
	// 			slice = bytes.Trim(strconv.AppendQuoteRune(slice, rune(vs.Uint())), "'")
	// 		default:
	// 			goto fool
	// 		}
	// 	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
	// 		// debug("------> 3 float")
	// 		switch f.verb {
	// 		case 'v':
	// 			if f.dot {
	// 				if f.prec == 0 {
	// 					slice = strconv.AppendFloat(slice, vs.Float(), 'e', 0, 64)
	// 				} else {
	// 					slice = strconv.AppendFloat(slice, vs.Float(), 'f', 0, 64)
	// 				}
	// 			} else {
	// 				slice = strconv.AppendFloat(slice, vs.Float(), 'f', f.prec, 64)
	// 			}
	// 		default:
	// 			slice = strconv.AppendFloat(slice, vs.Float(), f.verb, f.prec, 64)
	// 		}
	// 	case reflect.String:
	// 		// debug("------> 4 string")
	// 		v := a.(string)
	// 		// v := vs.String()
	// 		switch f.verb {
	// 		case 's', 'v':
	// 			slice = []byte(v)
	// 		case 'q':
	// 			slice = strconv.AppendQuote(slice, v)
	// 		case 'x', 'X':
	// 			tb := []byte(v)
	// 			for _, b := range tb {
	// 				slice = strconv.AppendInt(slice, int64(b), 16)
	// 			}
	// 		}
	// 	case reflect.Slice, reflect.Array:
	// 		if vs.Type().Elem().Kind() == reflect.Uint8 { // binary data
	// 			// debug("------> 5 []byte")
	// 			// slice = append(slice, vs.Bytes()...)
	// 			v := a.([]byte)
	// 			switch f.verb {
	// 			case 'x', 'X':
	// 				for _, b := range v {
	// 					slice = strconv.AppendInt(slice, int64(b), 16)
	// 				}
	// 			default:
	// 				slice = v
	// 			}
	// 		} else {
	// 			goto fool
	// 			// slice = append(slice, []byte(fmt.Sprintf("%"+string(f.verb), a))...)
	// 			// debug("------> 6 []%s", vs.Type().Elem().Kind().String())

	// 			// for i := 0; i < vs.Len(); i++ {
	// 			// 	slice = append(slice, dbg.toBytes(f, vs.Index(i))...)
	// 			// }
	// 		}

	// 	case reflect.Ptr:
	// 		goto fool
	// 	case reflect.Struct:
	// 		goto fool
	// 	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Invalid, reflect.Map, reflect.Uintptr, reflect.UnsafePointer:
	// 		fallthrough
	// 	default:
	// 		// b := bytes.NewBuffer(make([]byte, 0, 65536))
	// 		// binary.Write(b, binary.LittleEndian, a)
	// 		// bs := b.Bytes()
	// 		// slice = append(slice, []byte(vs.String())...)
	// 		// dbg.debug("\t[toBytes] \x1b[31m7 unknown type: [%v]--[%+v]--<%#v>--<%s>\x1b[0m", bs, a, vs, vs.String())
	// 		goto fool
	// 	}
	// 	if f.upper {
	// 		slice = bytes.ToUpper(slice)
	// 	}
	// 	// dbg.debug("\t[toBytes] <<<< \x1b[32m%s\x1b[0m", slice)
	// 	tds = append(tds, ts{fn: "toBytes", s: stt, e: time.Now()})
	// 	return slice

	// fool:
	// 	// dbg.debug("\t[toBytes] <<<< \x1b[32m%s\x1b[0m by fool", slice)
	// 	dbg.fool = true
	// 	slice = append(slice, []byte(fmt.Sprintf("%"+string(f.verb), a))...)
	// 	tds = append(tds, ts{fn: "toBytes", s: stt, e: time.Now()})
	// 	return slice
}
func fool(f flag, a interface{}) []byte {
	return []byte(fmt.Sprintf("%"+string(f.verb), a))
}

// func parseIndexI(idx string, args []interface{}) (int, error) {
// 	idx = strings.TrimSuffix(idx, "[")
// 	idx = strings.TrimSuffix(idx, "]")

// 	if i, e := strconv.Atoi(idx); e == nil && i > 0 && len(args) >= i {
// 		tar := args[i-1]
// 		switch tar.(type) {
// 		case int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
// 			debug("args[%d]=%#v, type:%#v", i-1, tar, tar.(interface{}))
// 			return tar.(int), nil
// 		case int:
// 			return tar.(int), nil
// 		}
// 	}

// 	return valErr, errIndex
// }

// ToString converts all types of value to string
func ToString(a ...interface{}) string {
	// reserved := make([]byte, 32)
	// binary.Write(b, binary.LittleEndian, reserved)
	// binary.Write(b, binary.LittleEndian, uint32(len(data)))
	// return string(ToBytes(a))
	slice := make([]byte, 0)
	for _, s := range a {
		slice = append(slice, ToBytes(s)...)
	}
	return string(slice)
}

// ToBytes converts all types of value to []byte
func ToBytes(a interface{}) []byte {
	// debug("[%v]  ", a)
	// b := bytes.NewBuffer(make([]byte, 65536))
	slice := []byte{}
	vs := reflect.ValueOf(a)
	switch vs.Kind() {
	case reflect.Bool:
		// if v := vs.Bool(); v {
		// 	slice = []byte("true")
		// }
		// slice = []byte("false")
		// debug("------> 0 bool")
		slice = strconv.AppendBool(slice, vs.Bool())
		// b.Write()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// debug("------> 1 int")
		slice = strconv.AppendInt(slice, vs.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// debug("------> 2 uint")
		slice = strconv.AppendUint(slice, vs.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		// debug("------> 3 float")
		slice = strconv.AppendFloat(slice, vs.Float(), 'f', -1, 64)
	case reflect.String:
		// debug("------> 4 string")
		// slice = append(slice, []byte(vs.String())...)
		slice = append(slice, []byte(a.(string))...)
	case reflect.Slice, reflect.Array:
		if vs.Type().Elem().Kind() == reflect.Uint8 { // binary data
			// debug("------> 5 []byte")
			// slice = append(slice, vs.Bytes()...)
			slice = append(slice, a.([]byte)...)
		} else {
			// debug("------> 6 []%s", vs.Type().Elem().Kind().String())
			for i := 0; i < vs.Len(); i++ {
				slice = append(slice, ToBytes(vs.Index(i))...)
			}
		}

	case reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Invalid, reflect.Map, reflect.Ptr, reflect.Struct, reflect.Uintptr, reflect.UnsafePointer:
		fallthrough
	default:
		// b := bytes.NewBuffer(make([]byte, 0, 65536))
		// binary.Write(b, binary.LittleEndian, a)
		// bs := b.Bytes()
		slice = append(slice, []byte(vs.String())...)
		// debug("------> 7 unknown type: %s --[%v]--[%+v]--<%#v>--<%s> ", vs.Kind().String(), bs, a, vs, vs.String())
	}
	return slice
}

// IntToBytes 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToInt 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func init() {
	// os.Remove("fmt.debug")
}

type dg struct {
	on, fool bool
	msg      []string
}

func (dbg *dg) debug(f string, a ...interface{}) {
	// return
	m := f
	if len(a) > 0 {
		m = fmt.Sprintf(f, a...)
	}
	if dbg.on {
		// fmt.Println(m)
	}
	dbg.msg = append(dbg.msg, m)
	// return

	fi, err := os.OpenFile("fmt.debug", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return
	}
	defer fi.Close()
	fi.WriteString(m + "\n")
	// return
}
func (dbg *dg) record() {
	fi, err := os.OpenFile("fmt.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return
	}
	defer fi.Close()

	for _, m := range dbg.msg {
		fi.WriteString(m + "\n")
	}
}
