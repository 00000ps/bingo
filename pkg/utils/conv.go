package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/jinzhu/copier"
)

const (
	// UTF8 charset name
	UTF8 = "utf-8"
	// GBK charset name
	GBK = "gbk"
)

// GBK2UTF8 convert string from gbk to utf-8
func GBK2UTF8(str string) string {
	return Any2UTF8(str, GBK)
}

// UTF82GBK convert string from utf-8 to gbk
func UTF82GBK(str string) string {
	return UTF82Any(str, GBK)
}

// Any2UTF8 convert string from specified charset to utf-8
func Any2UTF8(str string, charset string) string {
	dec := mahonia.NewDecoder(charset)
	return dec.ConvertString(str)
}

// UTF82Any convert string from utf-8 to specified charset
func UTF82Any(str string, charset string) string {
	enc := mahonia.NewEncoder(charset)
	return enc.ConvertString(str)
}

// Unicode2UTF8 will try to convert text in unicode to chinese charset
func Unicode2UTF8(text string) string {
	if !strings.Contains(text, "\\u") {
		if strings.Contains(text, "0026#") {
			text = strings.Replace(text, "0026#", "\\u0026#", -1)
			// log.Error("111   %s", text)
		} else if strings.Contains(text, "&") {
			text = strings.Replace(text, "&", "\\u0026", -1)
			// log.Error("2222   %s", text)
		} else {
			return text
		}
	}

	sUnicodev := strings.Split(text, "\\u")
	var context string
	for _, v := range sUnicodev {
		// fmt.Println(context)
		// fmt.Println(v)
		if len(v) < 1 {
			continue
		}
		if len(v) < 4 {
			context += v
			continue
		}
		var left string
		var temp int64
		var err error
		i := strings.Index(v, "#")
		// log.PrintRaw(i)
		// ts := ""
		if i == 4 {
			// context += strings.Join(vs[:len(vs)-1], "#")
			// ts := vs[len(vs)-1]
			ts := v[5:]
			// log.Notice(ts)
			j := strings.Index(ts, ";")
			if len(ts) < 2 || j == 0 {
				context += v
				continue
			}
			if j < 0 {
				temp, err = strconv.ParseInt(ts[:5], 0, 32)
				left = ts[6:]
			} else {
				temp, err = strconv.ParseInt(ts[:j], 0, 32)
				left = ts[j+1:]
			}
			// if len(ts) > 5 {
			// 	left = ts[6:]
			// }
		} else {
			temp, err = strconv.ParseInt(v[:4], 16, 32)
		}
		if err != nil {
			//log.Error("%s(%d)%s--%s", v, len(v), v[:4], err)
			context += v
			continue
		}
		// fmt.Println(v)
		// fmt.Printf("111   %#v\n", temp)
		// fmt.Printf("222   %c\n", temp)
		context += fmt.Sprintf("%c", temp)
		if len(v) > 4 {
			if i == 4 {
				context += left
			} else {
				context += v[4:]
			}
		}
	}
	return context
}

// UTF82Unicode will try to convert text in chinese charset to unicode
func UTF82Unicode(text string) string {
	textQuoted := strconv.QuoteToASCII(text)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	return textUnquoted
}

// Copy copier库就能处理不同类型之间的赋值。除此之外，copier还能：
// 调用同名方法为字段赋值；
// 以源对象字段为参数调用目标对象的方法，从而为目标对象赋值（当然也可以做其它的任何事情）；
// 将切片赋值给切片（可以是不同类型哦）；
// 将结构体追加到切片中
// https://darjun.github.io/2020/03/13/godailylib/copier/
func Copy(from, to interface{}) { copier.Copy(&to, &from) }

var ByteOrder = binary.LittleEndian

func Convert2FileName(name string) string {
	name = strings.Replace(name, ":", "_", -1)
	name = strings.Replace(name, " ", "_", -1)
	return strings.Replace(name, "__", "_", -1)
}

func Bytes2Int8(bytes []byte) int8 {
	return int8(bytes[0])
}
func Bytes2Int16(bytes []byte) int16 {
	return int16(ByteOrder.Uint16(bytes))
}
func Bytes2Int32(bytes []byte) int32 {
	return int32(ByteOrder.Uint32(bytes))
}
func Bytes2Int64(bytes []byte) int64 {
	return int64(ByteOrder.Uint64(bytes))
}
func Bytes2Int(bytes []byte) int {
	return int(ByteOrder.Uint32(bytes))
}

func Any2Bytes(data interface{}) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, ByteOrder, data)

	return b_buf.Bytes()
}
func Bytes2Any(data []byte) interface{} {
	var rd interface{}
	binary.Read(bytes.NewReader(data), ByteOrder, &rd)
	return rd
}
