package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

// JSONPretty returns formatted string of raw
func JSONPretty(data []byte) ([]byte, error) {
	js, err := simplejson.NewJson(data)
	if err != nil {
		return []byte{}, err
	}
	return js.EncodePretty()
}

// JSONBeauty returns data of v
func JSONBeauty(v interface{}) ([]byte, error) { return json.MarshalIndent(&v, "", "  ") }

// JSON returns string of v
func JSON(v interface{}) string {
	d, _ := json.MarshalIndent(&v, "", "  ")
	return string(d)
}

// JSONByte returns string of v
func JSONByte(v interface{}) []byte {
	d, _ := json.Marshal(&v)
	return d
}

// JSONPlain returns string of v
func JSONPlain(v interface{}) string {
	d, _ := json.Marshal(&v)
	return string(d)
}

// JSONFormat returns string of raw data
func JSONFormat(data []byte) string {
	njson := string(data)
	// njson = strings.Replace(njson, " ", "", -1)
	njson = strings.Replace(njson, "},{\"", "\n},\n{\n\"", -1)
	njson = strings.Replace(njson, "},\"", "\n},\n\"", -1)
	njson = strings.Replace(njson, "{\"", "{\n\"", -1)
	njson = strings.Replace(njson, ",\"", ",\n\"", -1)
	njsons := strings.Split(njson, "\n")

	for i, line := range njsons {
		line = strings.Replace(line, "\":", "\": ", 1)
		if !strings.Contains(line, `\"`) && len(line) > 2 {
			// line = strings.Replace(line, "{", "\n{", -1)
			line = strings.Replace(line, "}", "\n}", -1)
			// line = strings.Replace(line, "[", "\n[", -1)
			line = strings.Replace(line, "[", "[\n", -1)
			line = strings.Replace(line, "]", "\n]", -1)
		}
		if strings.HasSuffix(line, `"}`) && !strings.HasSuffix(line, `\"}`) {
			line = strings.TrimSuffix(line, `}`) + "\n}"
		}
		njsons[i] = line
	}

	njson = strings.Join(njsons, "\n")
	njsons = strings.Split(njson, "\n")
	level := 0
	for i, line := range njsons {
		njsons[i] = strings.Repeat(" ", level) + line
		if strings.HasPrefix(line, "{") || strings.HasSuffix(line, "{") || strings.HasPrefix(line, "[") || strings.HasSuffix(line, "[") {
			level++
			// log.Warning("add %d", level)
		} else if strings.HasSuffix(line, "}") || strings.HasSuffix(line, "},") || strings.HasSuffix(line, "]") || strings.HasSuffix(line, "],") {
			level--
			// log.Warning("sub %d", level)
			njsons[i] = strings.Repeat(" ", level*2) + line
		} else {
		}
		// fmt.Println("----------" + njsons[i])
	}
	njson = strings.Join(njsons, "\n")
	// njson = strings.Replace(njson, "", "", -1)
	return njson
}

// BeJSON identifies whether the input file is JSON format
func BeJSON(file string) bool {
	rawData, err := ReadFile(file)
	if err != nil {
		return false
	}

	var res map[string]interface{}
	err = json.Unmarshal(rawData, &res)
	if err != nil {
		var res1 []map[string]interface{}
		err = json.Unmarshal(rawData, &res1)
		if err != nil {
			return false
		}
	}
	return true
}

// FileSize returns the hommization file size
func FileSize(filename string) int64 {
	// var result int64
	fileInfo, _ := os.Stat(filename)
	return fileInfo.Size()

	// filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
	// 	result = f.Size()
	// 	return nil
	// })
	// return result
}

// FileSizeString returns the hommization file size
func FileSizeString(filename string) string {
	fileInfo, _ := os.Stat(filename)
	return Size(fileInfo.Size())
}

// Size returns the hommization file size
func Size(number interface{}) string {
	return FormatSize(number, 6)
}

// FormatSize returns the hommization file size
func FormatSize(number interface{}, prec int) string {
	var s float64
	switch number.(type) {
	case int:
		s = float64(number.(int))
	case int8:
		s = float64(number.(int8))
	case int16:
		s = float64(number.(int16))
	case int32:
		s = float64(number.(int32))
	case int64:
		s = float64(number.(int64))

	case uint:
		s = float64(number.(uint))
	case uint8:
		s = float64(number.(uint8))
	case uint16:
		s = float64(number.(uint16))
	case uint32:
		s = float64(number.(uint32))
	case uint64:
		s = float64(number.(uint64))

	case float32:
		s = float64(number.(float32))
	case float64:
		s = number.(float64)

	default:
		v := reflect.ValueOf(number)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			s = float64(v.Len())
		default:
			return "error number type, it should be number of list"
		}
	}

	switch {
	case s >= 1024*1024*1024:
		return strconv.FormatFloat(s/(1024*1024*1024), 'f', prec, 64) + "G"
	case s >= 1024*1024:
		return strconv.FormatFloat(s/(1024*1024), 'f', prec, 64) + "M"
	case s >= 1024:
		return strconv.FormatFloat(s/1024, 'f', prec, 64) + "K"
	default:
		return strconv.FormatFloat(s, 'f', -1, 64)
	}
}

// Int64 will convert a integer to int64
func Int64(number interface{}) int64 {
	switch number.(type) {
	case int:
		return int64(number.(int))
	case int8:
		return int64(number.(int8))
	case int16:
		return int64(number.(int16))
	case int32:
		return int64(number.(int32))
	case int64:
		return number.(int64)

	case uint:
		return int64(number.(uint))
	case uint8:
		return int64(number.(uint8))
	case uint16:
		return int64(number.(uint16))
	case uint32:
		return int64(number.(uint32))
	case uint64:
		return int64(number.(uint64))

	default:
		return 0
	}
}

// Float64 will convert a float to float64
func Float64(number interface{}) float64 {
	switch number.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return float64(Int64(number))
	case float32:
		return float64(number.(float32))
	case float64:
		return number.(float64)

	default:
		return 0
	}
}

// Percent returns 100*a/b %
func Percent(a, b interface{}) float64 {
	if bb := Float64(b); b != 0 {
		return 100 * Float64(a) / bb
	}
	return 0
}

// PercentString returns 100*a/b %
func PercentString(a, b interface{}, prec int) string {
	return strconv.FormatFloat(Percent(a, b), 'f', prec, 64) + "%"
}

// CompareAnyWithString will tell whether interface{} and string b are equal
func CompareAnyWithString(a interface{}, b string) bool {
	pa := reflect.ValueOf(a)
	switch pa.Kind() {
	case reflect.Int:
		typeVal := strconv.FormatInt(int64(a.(int)), 10)
		return assert.EqualValues(new(testing.T), typeVal, b)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		typeVal := strconv.FormatInt(a.(int64), 10)
		return assert.EqualValues(new(testing.T), typeVal, b)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		typeVal := strconv.FormatUint(a.(uint64), 10)
		return assert.EqualValues(new(testing.T), typeVal, b)
	case reflect.Float32, reflect.Float64:
		typeVal, _ := strconv.ParseFloat(b, 64)
		return assert.EqualValues(new(testing.T), a.(float64), typeVal)
	case reflect.Bool:
		typeVal, _ := strconv.ParseBool(b)
		return assert.EqualValues(new(testing.T), a.(bool), typeVal)
	case reflect.String:
		return assert.EqualValues(new(testing.T), a.(string), b)
	default:
		return false
	}
}

// PathFormat wi
func PathFormat(raw string) string { return strings.Replace(raw, ":/", "://", -1) }

// Quantile returns the specified value of an array
func Quantile(slice sort.Interface, accuracy float32) (position int) {
	if !sort.IsSorted(slice) {
		sort.Sort(slice)
	}
	pos := int(float32(slice.Len())*accuracy) - 1
	if pos < slice.Len()-1 {
		return pos
	}
	return slice.Len() - 1
	/*pos := int(float32(len(slice))*accuracy) - 1
	if pos+1 <= len(slice) {
		log.Info("%.2f%%\t%f ms", accuracy*100, slice[pos].cost*1000)
	}*/
}

// Split used to convert a list string to an array
func Split(raw string, suffix ...string) []string {
	raw = strings.Replace(raw, ",", ";", -1)
	arr := strings.Split(raw, ";")
	if len(suffix) > 0 {
		for i := 0; i < len(arr); i++ {
			if !strings.Contains(arr[i], suffix[0]) {
				arr[i] += suffix[0]
			}
		}
	}
	return arr
}

func Format4(format string, args ...interface{}) string   { return FormatLen(4, format, args...) }
func Format8(format string, args ...interface{}) string   { return FormatLen(8, format, args...) }
func Format10(format string, args ...interface{}) string  { return FormatLen(10, format, args...) }
func Format12(format string, args ...interface{}) string  { return FormatLen(12, format, args...) }
func Format16(format string, args ...interface{}) string  { return FormatLen(16, format, args...) }
func Format20(format string, args ...interface{}) string  { return FormatLen(20, format, args...) }
func Format200(format string, args ...interface{}) string { return FormatLen(200, format, args...) }
func Format500(format string, args ...interface{}) string { return FormatLen(500, format, args...) }
func Format(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}
func FormatLen(length uint, format string, args ...interface{}) string {
	text := Format(format, args...)
	t := ""
	if length > 2 {
		if uint(len(text)+1) > length {
			t = text[:(length-2)] + ". "
		} else {
			t = text + strings.Repeat(" ", int(length)-len(text))
		}
	}
	return t
}
func FormatMax(length uint, format string, args ...interface{}) string {
	text := Format(format, args...)
	t := ""
	if length > 4 {
		if uint(len(text)+1) > length {
			t = text[:(length-4)] + "... "
		} else {
			t = text
		}
	}
	return t
}

//SplitIPPort return the ip and port arr by split input format
func SplitIPPort(addrs string) (ip, port string) {
	s := strings.Split(addrs, ":")
	if len(s) == 2 {
		ip = s[0]
		port = s[1]
	}
	return
}

//GetIPPort return the ip and port arr by split input format
func GetIPPort(addrs ...string) (ips, ports []string) {
	for i := 0; i < len(addrs); i++ {
		ip, port := SplitIPPort(addrs[i])
		ips = append(ips, ip)
		ports = append(ports, port)
	}
	return
}

//FirstUpper used to convert string in first upper and the other lower format
func FirstUpper(str string) string {
	str = strings.ToLower(str)
	if len(str) > 1 {
		return strings.ToUpper(string(str[0])) + str[1:]
	}
	return strings.ToUpper(str)
}

//FirstUpperRaw used to convert string in first upper and the other not changed format
func FirstUpperRaw(str string) string {
	if len(str) == 0 {
		return str
	}
	fstr := string(str[0])
	ufstr := strings.ToUpper(fstr)
	return ufstr + str[1:]
}
