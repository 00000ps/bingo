package stub

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"bingo/internal/app/frame/testmonitor"
	"bingo/internal/app/frame/walle"
	"bingo/pkg/encode"
	"bingo/pkg/utils"
	"bingo/pkg/utils/format"
	"bingo/pkg/utils/log"
)

type (
	response struct {
		key      int64
		value    []byte
		duration time.Duration
	}
	Mock struct {
		//caselogger
		log.Logger
		Self       string
		Other      string
		Port       int
		fixedport  bool
		id         int
		timeout    uint16
		ln         net.Listener
		handler    func(net.Conn)
		Snapshot   func(d []byte, tag string)
		closer     bool
		indata     [][]byte
		outdata    []byte
		responses  map[int64]*response
		reqCounter int
		resCounter int
		bFile      string
	}
)

func (m *Mock) Count() int {
	if check(m) != nil {
		return -1
	}
	return m.reqCounter
}
func (m *Mock) BindResTimeout(timeoutMs uint16) {
	if check(m) != nil {
		return
	}
	m.timeout = timeoutMs
}
func (m *Mock) BindHandler(h func(net.Conn)) {
	if check(m) != nil {
		return
	}
	m.handler = h
}
func (m *Mock) Bind(file string, response ...interface{}) {
	if check(m) != nil {
		return
	}
	if err := m.bind(file, response...); err != nil {
		m.Error("[Mock:%s--%s] bind data error: %s", m.Self, m.Other, err)
	}
}
func (m *Mock) bind(file string, response ...interface{}) error {
	file = walle.GetPipeFile(file)
	data, err := utils.ReadFile(file)
	if err != nil {
		//m.Error("[Mock:%s--%s] bind failed: %s", m.Self, m.Other, err)
		return err
	}
	m.bFile = file
	m.bindData(data, response...)
	m.Notice("[Mock:%s--%s] load data: %s", m.Self, m.Other, format.Blue(m.bFile))
	return nil
}
func (m *Mock) bindData(data []byte, response ...interface{}) {
	m.outdata = data
	if len(response) > 0 {
		m.ParseResponse(response[0])
	}
}
func (m *Mock) BindResponse(response interface{}) []byte {
	if check(m) != nil {
		return []byte{}
	}
	ret, msg := encode.SmartEncoder(response)
	if msg.Error() != nil {
		m.Error("[Mock:%s->%s] bind failed: %s", m.Self, m.Other, msg.Error())
	}
	if msg.Warning() != nil {
		m.Warning("[Mock:%s->%s] encoding: %s", m.Self, m.Other, msg.Warning())
	}
	//m.Flash(data, fmt.Sprintf("Mock-%s-%s", m.Self, m.Other))
	m.bindData(ret.Data)
	return ret.Data
}

func (m *Mock) ParseRequest(request interface{}, index ...int) interface{} {
	if check(m) != nil {
		return nil
	}
	l := len(m.indata)
	if l == 0 {
		return nil
	}
	var data []byte
	if len(index) == 0 {
		data = m.indata[l-1]
	} else {
		if index[0]+1 >= l {
			return nil
		}
		data = m.indata[index[0]]
	}

	if err := encode.SmartDecoder(data, request); err != nil {
		m.Error("[Mock:%s--%s] parse request error: %s", m.Self, m.Other, err)
	}
	if request == nil {
		m.Error("[Mock:%s--%s] parse request error: request is nil", m.Self, m.Other)
	}
	return request
}
func (m *Mock) ParseResponse(response interface{}) {
	if check(m) != nil {
		return
	}
	if len(m.outdata) == 0 {
		return
	}

	if err := encode.SmartDecoder(m.outdata, response); err != nil {
		m.Error("[Mock:%s--%s] parse response error: %s", m.Self, m.Other, err)
	}
	if response == nil {
		m.Error("[Mock:%s--%s] parse response error: response is nil", m.Self, m.Other)
	}
	return
}

func NewMock(logger log.Logger, name, src string, tag, port int, fixedport bool, timeoutMs uint16, snapshot func([]byte, string)) (m *Mock) {
	m = new(Mock)
	m.Logger = log.SetLogger(logger)
	m.Self = name
	m.Other = src + strconv.Itoa(tag)
	m.bind(datapath(m.Self, src, false))
	if port == 0 {
		p, err := utils.GetPort()
		if err != nil {
			m.Error("[Mock:%s--%s] create mock failed due to set random port: %s", m.Self, m.Other, err)
			return
		}
		//m.Port = strconv.Itoa(p)
		m.Port = p
		m.Notice("[Mock:%s--%s] create mock on random port %d", m.Self, m.Other, m.Port)
	} else {
		//m.Port = strconv.Itoa(port)
		m.Port = port
		m.Notice("[Mock:%s--%s] create mock on specified port %d", m.Self, m.Other, m.Port)
	}
	m.fixedport = fixedport
	m.closer = false
	m.reqCounter = 0
	m.timeout = timeoutMs
	m.handler = nil
	m.id = testmonitor.Reg(m.Self, testmonitor.STUB, m, m.Off, m.Logger)
	m.Snapshot = snapshot
	//mockList[m.id] = m
	go m.listen()
	time.Sleep(time.Millisecond * 100)

	return m
}
func (m *Mock) Off() error {
	if m == nil {
		return log.NewError("err: mock is nil")
	}
	if m.reqCounter == 0 {
		m.Notice("[Mock:%s--%s] close mock, %s", m.Self, m.Other, format.Warning("did NOT receive any request"))
	} else {
		m.Notice("[Mock:%s--%s] close mock, request count: %d", m.Self, m.Other, m.reqCounter)
	}
	m.closer = true
	if m.ln != nil {
		m.ln.Close()
	}
	testmonitor.Unreg(m.id)
	//delete(mockList, m.id)
	if m != nil {
		m = nil
	}
	return nil
}
func (m *Mock) Close() {
	m.Off()
}

func (m *Mock) listen() {
	var err error
	m.ln, err = net.Listen("tcp", ":"+strconv.Itoa(m.Port))
	if err != nil {
		if m.fixedport {
			m.Error("[Mock:%s--%s] failed to listen port %d", m.Self, m.Other, m.Port)
			return
		} else {
			m.Notice("[Mock:%s--%s] listen port error, will try to reget a random port: %s", m.Self, m.Other, err)
			p, err := utils.GetPort()
			if err != nil {
				m.Error("[Mock:%s--%s] failed to get random port: %s", m.Self, m.Other, err)
				return
			}
			//m.Port = strconv.Itoa(p)
			m.Port = p
			m.ln, err = net.Listen("tcp", ":"+strconv.Itoa(m.Port))
			if err != nil {
				m.Error("[Mock:%s--%s] listen port error: %s", m.Self, m.Other, err)
				return
			}
		}
	}
	m.Notice("[Mock:%s--%s] listen port on %d", m.Self, m.Other, m.Port)
	for {
		var conn net.Conn
		conn, err = m.ln.Accept()
		if err != nil {
			if m.closer {
				return
			}
			m.Error("[Mock:%s<-%s] create connection error: %s", m.Self, m.Other, err)
			return
		}
		go func() {
			if m.closer {
				return
			}
			m.reqCounter++
			//m.Notice("[Mock:%s<-%s] %d: ready to receive request from %s", m.Self, m.Other, m.reqCounter, conn.RemoteAddr())

			defer conn.Close()
			var recvdata []byte
			recvdata, err = encode.ConnRead(conn)
			if err != nil {
				m.Error("[Mock:%s<-%s] %d: receive request len: %d, error: %s", m.Self, m.Other, m.reqCounter, len(recvdata), err)
				return
			}
			m.indata = append(m.indata, recvdata)
			m.Snapshot(recvdata, fmt.Sprintf("mock.REQ.%s-%s.%d", m.Other, m.Self, m.reqCounter))
			m.Notice("[Mock:%s<-%s] %d: received request len: %d", m.Self, m.Other, m.reqCounter, len(recvdata))

			if m.timeout > 0 {
				t := time.Millisecond * time.Duration(m.timeout)
				m.Notice("[Mock:%s--%s] %d: ready to execute timeout: %s", m.Self, m.Other, m.reqCounter, t)
				time.Sleep(t)
			}
			if m.handler != nil {
				m.Notice("[Mock:%s->%s] %d: call handler to %s", m.Self, m.Other, m.reqCounter, conn.RemoteAddr())
				m.handler(conn)
			} else {
				m.Snapshot(m.outdata, fmt.Sprintf("mock.RES.%s-%s.%d", m.Self, m.Other, m.reqCounter))
				wl, err := conn.Write(m.outdata)
				if err != nil {
					m.Error("[Mock:%s->%s] %d: send response len: %d, error: %s", m.Self, m.Other, m.reqCounter, wl, err)
					return
				}
				m.Notice("[Mock:%s->%s] %d: send response len: %d to %s", m.Self, m.Other, m.reqCounter, len(m.outdata), conn.RemoteAddr())
				if len(m.outdata) <= 0 {
					m.Warning("[Mock:%s->%s] %d: response len: %d, did you forget binding data for mock?", m.Self, m.Other, m.reqCounter, len(m.outdata))
				}
			}
		}()
	}
}
