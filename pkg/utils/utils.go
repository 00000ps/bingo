package utils

import (
	"bingo/pkg/cmd"
	"bingo/pkg/log"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"math"
	"net"
	"net/mail"
	"net/smtp"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/scorredoira/email"
	"github.com/stretchr/testify/assert"
)

func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

type HotConf struct {
	file         string
	conf         interface{}
	cycle        time.Duration
	getData      func() []byte
	unmarshaler  func([]byte, interface{}) error
	errorHandler func()
	md5sum       [md5.Size]byte
}

func HotLoad(data []byte, cycleSec int) []byte {
	// md5.Sum()
	// confdata := ""
	// for {
	// 	if fromIcafe {
	// 		cards, _ := icafe.GetCard(icafeSpace, icafeID, false)
	// 		if len(cards) == 1 {
	// 			card = cards[0]
	// 			confdata = strings.Replace(card.GetField("背景"), "&amp;quot;", "\"", -1)
	// 		}
	// 	} else {
	// 		if data, err := utils.ReadFile(file); err != nil {
	// 			log.Error("load conf error: %s: %s", err, file)
	// 			// file = defaultConf
	// 		} else {
	// 			confdata = string(data)
	// 		}
	// 	}
	// 	if m := utils.MD5(confdata); m != md5 {
	// 		if fromIcafe {
	// 			log.Notice("read conf from iCafe: %s", card.PubProperty.URL)
	// 		} else {
	// 			log.Notice("read conf from file: %s", file)
	// 		}
	// 		md5 = m
	// 		c = New([]byte(confdata))
	// 		if c != nil {
	// 			goCron()
	// 		}
	// 	}
	// 	time.Sleep(waitTime)
	// }
	// c := &conf{}
	// var err error
	// if err = json.Unmarshal(data, c); err != nil {
	// 	log.Error("parse conf error: %s", err)
	// 	if card != nil {
	// 		robot.Send("Conf Modified Error", `修改人: `+card.LastModifiedUser.Name+`\n修改时间: `+card.LastModifiedTime+`\n错误: `+err.Error(), global.HiFaceTeam)
	// 	}
	// 	return nil
	// }
	// if changeid > 0 {
	// 	log.Notice("conf changed: %d", changeid)
	// 	content := "修改次数: " + strconv.Itoa(int(changeid))
	// 	if card != nil {
	// 		content += `\n修改人员: ` + card.LastModifiedUser.Name + `\n修改时间: ` + card.LastModifiedTime
	// 	}
	// 	dm, err := diff.JSONCompare(data, confJson, false)
	// 	if err != nil {
	// 		log.Error("compare diff error: %s", err)
	// 		content += "\n获取修改内容错误: " + err.Error()
	// 	}
	// 	for k, v := range dm {
	// 		content += fmt.Sprintf("\n修改内容: %s: %v->%v", k, v.ControlValue, v.TestValue)
	// 	}
	// 	log.Warning(content)
	// 	if card != nil {
	// 		robot.Send("Conf Modified", content, global.HiFaceTeam)
	// 	}
	// }
	// confJson = data
	// changeid++
	// if c.Wait > 0 {
	// 	waitTime = time.Second * time.Duration(c.Wait)
	// }
	// for k, s := range c.Tasks {
	// 	s.key = k
	// 	if s.Enable {
	// 		s.Sched.Spec = strings.TrimSpace(s.Sched.Spec)
	// 		if s.Sched.Spec == "" {
	// 			log.Warning("task <%s> spec is nil, notification will be CLOSED", s.key)
	// 			// continue
	// 		}
	// 		s.ext.since, err = utils.ParseDate((utils.DateLayout, s.Sched.StartTime)
	// 		if err != nil {
	// 			log.Error("task <%s> parse start time error: %s, it should be %s", s.key, err, utils.DateLayout)
	// 			continue
	// 		}
	// 		s.ext.forever = false
	// 		if strings.TrimSpace(s.Sched.EndTime) == "" {
	// 			s.ext.forever = true
	// 		} else {
	// 			s.ext.end, err = utils.ParseDate((utils.DateLayout, s.Sched.EndTime)
	// 			if err != nil {
	// 				log.Error("task <%s> parse end time error: %s", s.key, err)
	// 				s.ext.forever = true
	// 			}
	// 		}
	// 		if s.ext.forever {
	// 			log.Warning("task <%s> in NON-STOP mode", s.key)
	// 		}
	// 		if len(s.Sched.Dutys) == 0 {
	// 			log.Warning("task <%s> duty list is nil", s.key)
	// 			// continue
	// 		} else {
	// 			if u, ok := GetUnit[strings.ToLower(s.Sched.Unit)]; !ok {
	// 				log.Warning("task <%s> invalid unit settings: %s, only day/week/month supported, set to week defaultly", s.key, s.Sched.Unit)
	// 				s.ext.u = unitWeek
	// 			} else {
	// 				s.ext.u = u
	// 			}
	// 		}
	// 		// s.ext.e = make(map[time.Time]string)
	// 		var dates []string
	// 		for tstr := range s.Sched.Exception {
	// 			dates = append(dates, tstr)
	// 		}
	// 		// log.PrintRaw(dates)
	// 		sort.Strings(dates)
	// 		// log.PrintRaw(dates)
	// 		for _, tstr := range dates {
	// 			if name, ok := s.Sched.Exception[tstr]; ok {
	// 				t, err := utils.ParseDate((utils.DateLayout, tstr)
	// 				if err != nil {
	// 					log.Warning("task <%s> parse start time error: %s, it should be %s, skipped: %s %s", s.key, err, utils.DateLayout, tstr, name)
	// 				} else {
	// 					s.ext.e = append(s.ext.e, es{t, name})
	// 				}
	// 			}
	// 		}
	// 	}
	// 	d, _ := json.Marshal(s)
	// 	rawConfMap[k] = string(d)
	// 	log.Debug("task <%s>: %s", s.key, rawConfMap[s.key])
	// }
	// return c
	return []byte{}
}

func SendMail(subject, to, body string, htmlmail bool, files ...string) error {
	host := "proxy-in..com:25"
	hp := strings.Split(host, ":")
	//smtp.SendMail("addr string", smtp.PlainAuth("", user, "", hp[0]), "from string", "to []string", "msg []byte")
	user := "bingo"
	userAddr := user + "@.com"
	password := ""

	to = strings.Replace(to, ",", ";", -1)
	if gitUser := cmd.GetCurrentUserByGit(); gitUser != "" {
		to += ";" + gitUser
	}
	sendTo := strings.Split(to, ";")
	for i := 0; i < len(sendTo); i++ {
		if !strings.Contains(sendTo[i], "@") {
			sendTo[i] += "@.com"
		}
	}
	//msg := []byte("To: " + strings.Join(sendTo, ";") + "\r\nFrom: " + userAdd + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	//log.Notice("Send mail")
	//return smtp.SendMail(host, smtp.PlainAuth("", userAdd, password, hp[0]), user, sendTo, msg)
	var m *email.Message
	//contentType := "Content-Type: text/html" + "; charset=UTF-8"
	subject = "[bingo]" + subject
	if htmlmail {
		m = email.NewHTMLMessage(subject, body)
		//contentType = "Content-Type: text/html; charset=UTF-8"
		//} else {
		//contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	} else {
		m = email.NewMessage(subject, body)
	}
	//m := email.NewMessage(subject, contentType+"\r\n\r\n"+body)
	m.From = mail.Address{Name: user, Address: userAddr}
	m.To = sendTo
	for _, f := range files {
		err := m.Attach(f)
		if err != nil {
			log.Error("%s", err)
		}
	}
	return email.Send(host, smtp.PlainAuth("", userAddr, password, hp[0]), m)
}
func SendMailFile(subject, to, file string, htmlmail bool) error {
	body, _ := ReadFile(file)
	return SendMail(subject, to, string(body), htmlmail)
}

func GetEmail(name string) string { return strings.ToLower(name) + `@.com` }

func ParseCmd(cmd string) (k, v string) {
	if strings.HasPrefix(cmd, "--") {
		if strings.Contains(cmd, "=") {
			i := strings.Index(cmd, "=")
			return cmd[2:i], cmd[i+1:]
		}
		// return strings.TrimPrefix(cmd, "--"), ""
	}
	return "", ""
}

type Param struct{ Key, Value string }

func ParseFlag(cmd []string) []Param {
	//[-d 22206 -para 2 -bench 10]
	kv := make(map[string]string)
	var keys []string

	for i, c := range cmd {
		//fmt.Println(i, c)
		if k, v := ParseCmd(c); k != "" {
			kv[k] = v
			keys = append(keys, k)
		} else if strings.HasPrefix(c, "-") {
			k := strings.TrimPrefix(c, "-")
			k = strings.TrimPrefix(k, "-")
			kv[k] = ""
			if !(len(cmd) > i+1 && !strings.HasPrefix(cmd[i+1], "-")) {
				keys = append(keys, k)
			}
		} else if i > 0 && strings.HasPrefix(cmd[i-1], "-") && !strings.Contains(cmd[i-1], "=") {
			k := strings.TrimPrefix(cmd[i-1], "-")
			k = strings.TrimPrefix(k, "-")
			kv[k] = c
			keys = append(keys, k)
		} else {
			kv[c] = ""
			keys = append(keys, c)
		}
		// fmt.Println("keys: ", keys)
		// fmt.Println("kv: ", kv)
	}
	var p []Param
	for _, k := range keys {
		p = append(p, Param{k, kv[k]})
	}
	return p
}

//GetIntLen returns length of int value
func GetIntLen(num int) int { return int(math.Log10(float64(num))) + 1 }

// IsSelf returns whether host is current running machine
func IsSelf(host string) bool { return cmd.GetIP() == cmd.GetIP(host) }

// IsDiff returns
func IsDiff(file1, file2 string) bool { return MD5(file1) == MD5(file2) }

func ParsePara(paraValue string) []int {
	var paraValueListSlice []int
	paraValueList := strings.Split(paraValue, ",")
	for i := 0; i < len(paraValueList); i++ {
		if intValue, err := strconv.Atoi(paraValueList[i]); err == nil {
			paraValueListSlice = append(paraValueListSlice, intValue)
		} else {
			log.Error("Para init error,recheck input? Auto setting para to 1")
			return []int{1}
		}
	}
	return paraValueListSlice
}

// GetRandom return a random int value by time seed
func GetRandom(scare int) int {
	if scare == 0 {
		return int(time.Now().UnixNano())
	}
	return int(time.Now().UnixNano()) % scare
}

func URLEncode(r string) string          { return url.QueryEscape(r) }
func URLDecode(r string) (string, error) { return url.QueryUnescape(r) }

// MD5 returns the md5 value of file
func MD5File(fileName string) string {
	content, inerr := ReadFileAll(fileName)
	if inerr == nil {
		md5h := md5.New()
		//io.Copy(md5h, content)
		md5h.Write(content)
		cipherStr := md5h.Sum(nil)
		return hex.EncodeToString(cipherStr)
		//return fmt.Sprintf("%x", md5h.Sum([]byte(""))) //md5
	}
	return ""
}

// MD5 returns the md5 value of content
func MD5(content string) string {
	md5h := md5.New()
	md5h.Write([]byte(content))
	cipherStr := md5h.Sum(nil)
	return hex.EncodeToString(cipherStr)
	//return fmt.Sprintf("%x", md5h.Sum([]byte(""))) //md5
}
func FileBase64(filename string) string {
	if data, err := ReadFileAll(filename); err == nil {
		//en_img = base64.b64encode(data)
		return base64.StdEncoding.EncodeToString(data)
	}
	return ""
}
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
func DecodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
func URLEncodeBase64(url []byte) string {
	return base64.URLEncoding.EncodeToString(url)
}
func URLDecodeBase64(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}
func ImageBase64(imgs ...string) string {
	images := ""
	//imgs := strings.Split(imgstr, ",")
	for _, img := range imgs {
		if IsFileExist(img) {
			images += FileBase64(img)
			images += ","
		}
	}
	return strings.TrimSuffix(images, ",")
}

// GetID returns timestamp as id in int64 format
func GetID() int64 { return time.Now().UnixNano() }

// GetIDStr returns timestamp as id in string format
func GetIDStr() string { return strconv.FormatInt(GetID(), 10) }

// GetUUID returns uuid
func GetUUID() []byte { return uuid.NewUUID() }

// GetUUIDInt returns uuid in int64 format
func GetUUIDInt() int64 { return int64(binary.LittleEndian.Uint64(GetUUID())) }

// GetUUIDStr returns uuid in int64 format
func GetUUIDStr() string { return string(GetUUID()) }

// GetUUIDByLen returns uuid by speficied length in string format
func GetUUIDByLen(length int) string {
	time.Sleep(time.Nanosecond)
	m := MD5(strconv.FormatInt(GetID(), 10))
	if len(m) > length && length > 0 {
		return m[:length]
	}
	return m
}

// GetTimeIDStr returns id in string format
func GetTimeIDStr() string {
	return time.Now().Format("20060102150405")
}

// GetTimeID returns id in int64 format
func GetTimeID() int64 {
	/*
		#string到int
		int,err:=strconv.Atoi(string)
		#string到int64
		int64, err := strconv.ParseInt(string, 10, 64)
		#int到string
		string:=strconv.Itoa(int)
		#int64到string
		string:=strconv.FormatInt(int64,10)
	*/
	if i, err := strconv.ParseInt(time.Now().Format("20060102150405000"), 10, 64); err != nil {
		return GetID()
	} else {
		return i
	}
}
func GetPort() (int, error) {
	//rand.Seed(time.Now().UnixNano())
	//return strconv.Itoa(7000 + rand.Intn(8000-7000))
	// 得到一个可用的端口.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	//fmt.Println(listener.Addr())
	addr := listener.Addr().String()
	_, portString, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(portString)
}
func GetPWD() string { return os.Getenv("PWD") }

func Fuzzy(format string, input []string) (output []string) {
	sep := "*"
	if c := strings.Count(format, sep); c == 0 {
		for _, v := range input {
			if v == format {
				output = append(output, v)
			}
		}
		return
	} else if c == 1 {
		if strings.HasPrefix(format, sep) {
			for _, v := range input {
				if strings.HasSuffix(v, strings.TrimPrefix(format, sep)) {
					output = append(output, v)
				}
			}
		} else if strings.HasSuffix(format, sep) {
			for _, v := range input {
				if strings.HasPrefix(v, strings.TrimSuffix(format, sep)) {
					output = append(output, v)
				}
			}
		} else {
			cs := strings.Split(format, sep)
			if len(cs) > 2 {
				for _, v := range input {
					if strings.HasPrefix(v, cs[0]) && strings.HasSuffix(v, cs[1]) {
						output = append(output, v)
					}
				}
			}
		}
		return
	} else {
		return input
	}
}

func Assert(a, b interface{}) bool {
	return assert.Equal(&testing.T{}, a, b, "")
}

type Diff struct {
	standard, test interface{}
}

func NewDiff(standard interface{}) *Diff {
	d := new(Diff)
	d.standard = standard
	return d
}
func (d *Diff) Assert(test interface{}) bool {
	return Assert(d.standard, test)
}
