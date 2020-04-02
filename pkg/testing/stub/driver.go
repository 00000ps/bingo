package stub

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"bingo/internal/app/frame/testmonitor"
	"bingo/pkg/encode"
	"bingo/pkg/utils"
	"bingo/pkg/format"
	"bingo/pkg/utils"
	"bingo/pkg/log"
)

type (
	Driver struct {
		//caselogger
		log.Logger
		addr        url.URL
		rawaddr     string
		Conn        *net.Conn
		Snapshot    func([]byte, string)
		self, other string
		id          int
		resCounter  int
		// reqCounter, resCounter int
		counterChan chan int
		out         *encode.Result
		reqTimeout  time.Duration
		// HttpStatusCode int
		// buffer         map[int]buf
		binds struct {
			data []byte
			file string
			req  interface{}
		}
		diff *utils.Diff
		//protocal string
		//isparsed   bool
	}
	buf struct {
		code int
		in   []byte
		out  *encode.Result
	}
)

func (d *Driver) Bind(file string, request ...interface{}) {
	if check(d) != nil {
		return
	}
	if err := d.bind(file, request...); err != nil {
		d.Error("[Driver:%s--%s] bind data error: %s", d.self, d.other, err)
	}
}
func (d *Driver) bind(file string, request ...interface{}) error {
	file = walle.GetPipeFile(file)
	data, err := utils.ReadFile(file)
	if err != nil {
		return err
	}
	d.binds.file = file
	d.bindData(data, request...)
	d.Notice("[Driver:%s--%s] load data: %s", d.self, d.other, format.Blue(d.binds.file))
	return nil
}
func (d *Driver) bindData(data []byte, request ...interface{}) {
	d.binds.data = data
	d.out.Data = data
	if len(request) > 0 {
		d.binds.req = request
		d.ParseRequest(request[0])
	}
}

func (d *Driver) SetHost(host string) {
	if check(d) != nil {
		return
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = d.addr.Scheme + "://" + host
	}
	u, err := url.Parse(host + d.addr.Path + "?" + d.addr.RawQuery)
	if err != nil {
		d.Warning("Set driver host error: %s", err)
	}
	d.addr = *u
}
func (d *Driver) SetPath(p string) {
	if check(d) != nil {
		return
	}
	p = strings.TrimPrefix(p, "/")

	u, _ := url.Parse(d.rawaddr)

	raw := u.Scheme + "://" + u.Host + "/" + u.Path + "/" + p
	if !strings.Contains(p, "?") {
		raw += "?" + u.RawQuery
	}
	raw = path.Clean(raw)
	raw = utils.PathFormat(raw)
	nu, err := url.Parse(raw)
	if err != nil {
		d.Warning("Set driver path error: %s", err)
	}
	d.addr = *nu
}
func (d *Driver) SetQuery(query string) {
	if check(d) != nil {
		return
	}
	query = strings.TrimPrefix(query, "?")
	u, err := url.Parse(d.addr.Scheme + "://" + d.addr.Host + d.addr.Path + "?" + query)
	if err != nil {
		d.Warning("Set driver query error: %s", err)
	}
	d.addr = *u
}

func (d *Driver) SetStepAddr(addr string) {
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}
	d.rawaddr = addr
}
func (d *Driver) SetTimeout(to time.Duration) { d.reqTimeout = to }
func (d *Driver) SendTo(method string, request interface{}, succRes, errRes interface{}) (succ bool, res []byte, cost time.Duration, code int) {
	if check(d) != nil {
		return
	}
	d.SetPath(method)
	return d.Send(request, succRes, errRes)
}
func (d *Driver) Send(request interface{}, succRes, errRes interface{}) (succ bool, res []byte, cost time.Duration, code int) {
	//if !d.isparsed {
	//d.Error("[Driver:%s->%s] request did NOT been parsed, please call ParseResponse(request) first", d.self, d.other)
	//return false
	//}
	if check(d) != nil {
		return
	}
	d.binds.req = request

	ret, msg := encode.SmartEncoder(request)
	d.out = ret
	if msg.Error() != nil {
		d.Error("[Driver:%s->%s] encode failed: %s", d.self, d.other, msg.Error())
		return
	}
	if msg.Warning() != nil {
		d.Warning("[Driver:%s->%s] encoding: %s", d.self, d.other, msg.Warning())
	}
	if len(ret.Data) > 0 {
		a := mcpack.Decode(d.binds.data)
		b := mcpack.Decode(ret.Data)
		if !utils.Assert(a, b) {
			ts := "log/" + utils.GetTimestamp()
			rfile := ts + "_" + d.self + "-" + d.other + "_raw"
			nfile := ts + "_" + d.self + "-" + d.other + "_new"
			d.Warning("[Driver:%s->%s] encode warning! interface mismatched, please check: [vimdiff %s %s]", d.self, d.other, rfile, nfile)
			utils.WriteNewFile(rfile, []byte(a))
			utils.WriteNewFile(nfile, []byte(b))
			//d.Warning("[Driver:%s->%s] encode warning: maybe interface between %s & %s updated, please check", d.self, d.other, d.self, d.other)
			//fmt.Printf(color.Green("interface RAW: ========\n%s\n", a))
			//fmt.Printf(color.Red("interface NEW: ========\n%s\n", b))
		}
	} else if len(ret.Values) > 0 || ret.Buffer != nil {
		if !strings.HasPrefix(d.addr.String(), "http://") && !strings.HasPrefix(d.addr.String(), "https://") {
			d.Error("[Driver:%s->%s] incorrect url: %s", d.self, d.other, d.addr.String())
			return
		}
	}

	res, cost, code = d.send(ret)

	var errSucc, errFail error
	// res := d.GetResponse()
	errSucc = parseData(res, succRes)
	errFail = parseData(res, errRes)
	if errRes != nil {
		errFail = parseData(res, errRes)
	} else {
		if errSucc == nil {
			d.Info("[Driver:%s<-%s] success response", d.self, d.other)
			return
		}
	}
	//log.Notice("stub: parseData %s -- %s, %#v == %#v", errSucc, errFail, errRes, succRes)
	if errFail != nil && errSucc == nil {
		d.Info("[Driver:%s<-%s] success response", d.self, d.other)
		return
	}
	/*
		if errFail == nil && errSucc != nil {
			d.Notice("[Driver:%s<-%s] receiving success response", d.self, d.other)
		}
		if errSucc != nil && errFail != nil {
			d.Error("[Driver:%s<-%s] parse response error: %s %s", d.self, d.other, errSucc, errFail)
			return false
		}

		if errFail == nil {
			d.Notice("[Driver:%s<-%s] receiving success response", d.self, d.other)
			//return false
		}
	*/
	//if len(d.in.data) > 0 {
	//	d.Warning("[Driver:%s<-%s] error response", d.self, d.other)
	//}
	return
}
func (d *Driver) SendData(data []byte) (res []byte, cost time.Duration, code int) {
	if check(d) != nil {
		return
	}
	d.out.Data = data
	return d.send(d.out)
}
func (d *Driver) PostData(data url.Values) (res []byte, cost time.Duration, code int) {
	if check(d) != nil {
		return
	}
	d.out.Values = data
	return d.send(d.out)
}
func (d *Driver) SendDefaultData() (res []byte, cost time.Duration, code int) {
	if check(d) != nil {
		return
	}
	return d.send(d.out)
}
func (d *Driver) Count() int {
	if check(d) != nil {
		return -1
	}
	return d.resCounter
}
func (d *Driver) GetRequest() []byte {
	if check(d) != nil {
		return []byte{}
	}
	return d.out.Data
}

// func (d *Driver) GetResponse() []byte {
// 	if check(d) != nil {
// 		return []byte{}
// 	}
// 	if len(d.in.data) == 0 && d.out.Protocal != encode.ProtocalGET && d.out.Protocal != encode.ProtocalPOST {
// 		d.read()
// 	}
// 	return d.in.data
// }

func (d *Driver) ParseRequest(request interface{}) {
	if check(d) != nil {
		return
	}
	//d.isparsed = true
	if len(d.out.Data) == 0 {
		if file := datapath(d.self, d.other, true); utils.IsFileExist(file) {
			data, err := utils.ReadFile(file)
			if err != nil {
				d.Error("[Driver:%s--%s] parse request failed: %s", d.self, d.other, err)
				return
			}
			d.out.Data = data
		}
	}

	if err := parseData(d.out.Data, request); err != nil {
		d.Error("[Driver:%s--%s] parse request error: %s", d.self, d.other, err)
	}
	return
}

// func (d *Driver) ParseResponse(response interface{}) ([]byte, error) {
// 	if err := check(d); err != nil {
// 		return []byte{}, err
// 	}
// 	if err := parseData(d.GetResponse(), response); err != nil {
// 		return []byte{}, fmt.Errorf("[Driver:%s<-%s] parse response error: %s", d.self, d.other, err)
// 	}
// 	return d.in.data, nil
// }

func NewDrive(logger log.Logger, name, dst string, tag int, addr string, snapshot func([]byte, string)) (d *Driver) {
	d = new(Driver)
	d.self = name
	d.Logger = log.SetLogger(logger)
	d.other = dst + strconv.Itoa(tag)
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}
	d.rawaddr = addr
	u, err := url.Parse(addr)
	if err != nil {
		d.Warning("Set driver addr error: %s", err)
	}
	d.addr = *u
	d.out = new(encode.Result)
	// d.in.data = []byte{}
	d.resCounter = 0
	d.counterChan = make(chan int, 10000)
	go func() {
		i := 0
		for {
			d.counterChan <- i
			i++
		}
	}()
	// d.buffer = make(map[int]buf)
	d.bind(datapath(d.self, dst, true))
	d.id = testmonitor.Reg(d.self, testmonitor.STUB, d, d.Off, d.Logger)
	d.Snapshot = snapshot
	return
}
func (d *Driver) Off() error {
	if d == nil {
		return log.NewError("err: driver is nil")
	}
	if d.resCounter == 0 {
		d.Notice("[Driver:%s--%s] close driver, %s", d.self, d.other, format.Warning("did NOT receive any response"))
	} else {
		d.Notice("[Driver:%s--%s] close driver, response count: %d", d.self, d.other, d.resCounter)
	}
	if d.Conn != nil {
		(*d.Conn).Close()
		d.Conn = nil
	}
	testmonitor.Unreg(d.id)
	//delete(driverList, d.id)
	if d != nil {
		d = nil
	}
	return nil
}
func (d *Driver) Close() { d.Off() }

func (d *Driver) send(data *encode.Result) (res []byte, cost time.Duration, code int) {
	code = -1
	if len(data.Values) <= 0 &&
		len(data.EncodingData) <= 0 &&
		len(data.Data) == 0 &&
		data.Buffer == nil {
		d.Warning("[Driver:%s->%s] load data error, pls make sure %s exists or use Bind() to load your data", d.self, d.other, walle.GetReqPipeFile(d.self, d.other))
		//return
	}
	startTime := time.Now()
	// d.reqCounter++
	reqCounter := d.getReqCounter()
	// var c int
	switch data.Protocal {
	case encode.ProtocalGET, encode.ProtocalPOST:
		res, cost, code = d.sendHTTP(startTime, data, reqCounter)
	default:
		res, cost = d.sendTCP(startTime, data, reqCounter)
	}
	// d.buffer[reqCounter] = buf{
	// 	out:  data,
	// 	in:   res,
	// 	code: c,
	// }
	d.resCounter++
	if len(res) == 0 {
		d.Warning("[Driver:%s<-%s] reading response error: no data received, cost: %s", d.self, d.other, utils.GetCostStr(startTime))
	} else {
		d.Snapshot(res, fmt.Sprintf("driver.RES.%s-%s.%d", d.other, d.self, reqCounter))
	}
	return
}
func (d *Driver) sendTCP(startTime time.Time, data *encode.Result, reqCounter int) (res []byte, cost time.Duration) {
	start := time.Now()
	conn, err := net.Dial("tcp", d.addr.String())
	if err != nil {
		d.Error("[Driver:%s->%s] dial error: %s", d.self, d.other, err)
		return
	}
	d.Debug("[Driver:%s->%s] connecting cost: %s", d.self, d.other, utils.GetCostStr(startTime))

	d.Conn = &conn
	if d.reqTimeout > 0 {
		conn.SetDeadline(time.Now().Add(d.reqTimeout))
	}
	d.Debug("[Driver:%s->%s] ready to send request to %s", d.self, d.other, (*d.Conn).RemoteAddr())

	//startTime = time.Now()
	wl, err := conn.Write(data.Data)
	cost = time.Since(start)
	d.Notice("[Driver:%s->%s] send request len: %d, cost: %s", d.self, d.other, len(data.Data), utils.GetCostStr(startTime))
	d.Snapshot(data.Data, fmt.Sprintf("driver.REQ.%s-%s.%d", d.self, d.other, reqCounter))
	startTime = time.Now()
	if err != nil {
		d.Error("[Driver:%s->%s] send request len: %d, error: %s", d.self, d.other, wl, err)
		return
	}

	//d.read()
	d.Notice("[Driver:%s<-%s] trying to receive response", d.self, d.other)
	if d.Conn == nil {
		d.Error("[Driver:%s<-%s] receiving response error: disconnected", d.self, d.other)
		return
	}
	defer (*d.Conn).Close()
	res, err = encode.ConnRead(*d.Conn)
	if err != nil {
		d.Error("[Driver:%s<-%s] reading response len: %d, error: %s", d.self, d.other, len(res), err)
		return
	}
	return
}
func (d *Driver) sendHTTP(startTime time.Time, data *encode.Result, reqCounter int) (res []byte, cost time.Duration, code int) {
	var sendBody string
	var req *http.Request
	var err error
	var body io.Reader
	prefix := fmt.Sprintf("[Driver:%s->%s] %s", d.self, d.other, strings.ToUpper(data.Protocal))

	switch data.Encoding {
	case encode.EncodingFormData:
		body = data.Buffer
		sendBody = string(data.Buffer.Bytes())
	case encode.EncodingPROTOBUF:
		body = data.Buffer
		sendBody = string(data.Buffer.Bytes())
	case encode.EncodingURLEncoded:
		fallthrough
	default:
		if len(data.Values.Encode()) > 0 {
			sendBody = data.Values.Encode()
		} else if len(data.EncodingData) > 0 {
			sendBody = string(data.EncodingData)
		}
		body = strings.NewReader(sendBody)
	}

	req, err = http.NewRequest(strings.ToUpper(data.Protocal), d.addr.String(), body)

	if data.Header == nil || len(data.Header) <= 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header = data.Header
	}
	d.Notice("%s host: %s", prefix, format.Blue(d.addr.String()))
	d.Notice("%s header: %v", prefix, req.Header)
	start := time.Now()
	var client = http.DefaultClient
	if d.reqTimeout > 0 {
		client.Timeout = d.reqTimeout
	}
	resp, err := client.Do(req)

	cost = time.Since(start)
	if err != nil {
		d.Error("%s error, url=%s [%s], cost: %s", prefix, d.addr.String(), err, utils.GetCostStr(startTime))
		return
	}
	// d.HttpStatusCode = resp.StatusCode
	code = resp.StatusCode

	d.Notice("%s body: %s", prefix, format.Blue(utils.FormatMax(100, sendBody)))
	d.Notice("%s body. len: %d, cost: %s", prefix, len(sendBody), utils.GetCostStr(startTime))
	d.Snapshot([]byte(sendBody), fmt.Sprintf("driver.REQ.%s-%s.%d", d.self, d.other, reqCounter))
	d.Notice("[Driver:%s<-%s] trying to receive response", d.self, d.other)

	startTime = time.Now()
	defer resp.Body.Close()
	res, err = utils.Read(resp.Body)
	if err != nil {
		d.Error("[Driver:%s->%s] receiving response error, %s, cost: %s", d.self, d.other, err, utils.GetCostStr(startTime))
		return
	}
	// d.in.data = []byte(iconv.Unicode2UTF8(string(res)))
	d.Notice("[Driver:%s->%s] receiving response: %s", d.self, d.other, format.Blue(utils.FormatMax(100, iconv.Unicode2UTF8(string(res)))))

	//pr.Raw = string(res)
	//pr.Msg = iconv.Unicode2UTF8(pr.Raw)
	return
}
func (d *Driver) getReqCounter() int { return <-d.counterChan }

func (d *Driver) read() {
	d.Notice("[Driver:%s<-%s] trying to receive response", d.self, d.other)
	if d.Conn == nil {
		d.Error("[Driver:%s<-%s] receiving response error: disconnected", d.self, d.other)
		return
	}

	defer (*d.Conn).Close()
	data, err := encode.ConnRead(*d.Conn)
	if err != nil {
		d.Error("[Driver:%s<-%s] reading response len: %d, error: %s", d.self, d.other, len(data), err)
		return
	}
	d.resCounter++
	d.Snapshot(data, fmt.Sprintf("driver.RES.%s-%s.%d", d.other, d.self, d.resCounter))
	// d.in.data = data
	if len(data) == 0 {
		d.Warning("[Driver:%s<-%s] reading response error: no data received", d.self, d.other)
	} else {
		d.Notice("[Driver:%s<-%s] reading response len: %d from %s", d.self, d.other, len(data), (*d.Conn).RemoteAddr())
	}
}
