package encode

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"bingo/pkg/utils"
	"bingo/pkg/log"
)

const (
	HEADER_UNKNOWN = iota + 10
	HEADER_NSHEAD
	HEADER_CSHEAD
	HEADER_INT
	HEADER_INT16

	PROTOCAL_MCPACK = iota + 20
	PROTOCAL_PROTO
)

func ConnRead(conn net.Conn, way ...int) (data []byte, err error) {
	//defer conn.Close()
	use := 3
	if len(way) > 0 {
		use = way[0]
	}

	//conn.SetReadDeadline(time.Millisecond * 100)
	switch use {
	case 0:
		data, err = utils.Read(conn)
		if err != nil {
			return
		}
	case 1:
		return bufio.NewReader(conn).ReadBytes('\n')
	case 2:
		ret := bytes.NewBuffer(nil)
		var buf [1024]byte
		length := 0
		for {
			length, err = conn.Read(buf[0:])
			ret.Write(buf[0:length])
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}
		}
		data = ret.Bytes()
	default:
		conn.SetReadDeadline(time.Now().Add(time.Second * 30))

		read := func(expectLen, bakLen int) error {
			for {
				var l int
				buf := make([]byte, 1024)
				l, err = conn.Read(buf)
				data = append(data, buf[:l]...)
				//log.Info("l=%d; length=%d; datalen(%d)", l, length, len(data))
				if len(data) == expectLen {
					break
				} else if bakLen != 0 && len(data) == bakLen {
					break
				}
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
			}
			return nil
		}

		length := 0
		buffer := make([]byte, nshead.NSHEAD_LENGTH)
		length, err = conn.Read(buffer)
		if err != nil {
			return
		}
		data = buffer[:length]
		if nh, ok := nshead.Unmarshal(buffer); ok {
			//log.Info("len=%d, %#v", length, nh)
			//if length == nshead.NSHEAD_LENGTH {
			//	ok = true
			//}
			return data, read(int(nh.BodyLen)+nshead.NSHEAD_LENGTH, 0)
		} else if ch, ok := cshead.Unmarshal(buffer); ok {
			return data, read(int(ch.BodyLen)+cshead.CSHEAD_LENGTH, 0)
		}

		leng := int(binary.LittleEndian.Uint32(buffer[:4])) + 4
		baklen := int(binary.LittleEndian.Uint16(buffer[:2])) + 4
		if length >= 36 {
			if int64(binary.LittleEndian.Uint64(buffer[:8])) == 0 &&
				int64(binary.LittleEndian.Uint64(buffer[8:16])) == 0 &&
				int64(binary.LittleEndian.Uint64(buffer[16:24])) == 0 &&
				int64(binary.LittleEndian.Uint64(buffer[24:32])) == 0 {
				baklen = int(binary.LittleEndian.Uint32(buffer[32:36])) + 36
			}
		}
		return data, read(leng, baklen)

		/*var body []byte
		body, err = bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return
		}
		data = make([]byte, len(body))
		data = append(buffer, body...)*/
	}
	return
}

// decode data with specified header & protocal
func SpecifyParse(data []byte, v interface{}, header, protocal int) (length int, nh nshead.NsHead, err error) {
	var body []byte
	if header == HEADER_NSHEAD {
		n, _ := nshead.Unmarshal(data)
		nh = *n
		body = data[nshead.NSHEAD_LENGTH:]
		length = int(nh.BodyLen)
	} else if header == HEADER_INT {
		body = data[4:]
		length = int(conv.Bytes2Int32(data[:4]))
	} else if header == HEADER_INT16 {
		body = data[2:]
	} else {
		log.Err(UnsupportedHeaderErr)
	}

	if protocal == PROTOCAL_MCPACK {
		err = mcpack.Unmarshal(body, v)
		if err != nil {
			fmt.Printf("%#v\n", err)
			return
		}
	} else if protocal == PROTOCAL_PROTO {
		err = proto.Unmarshal(body, v.(proto.Message))
		if err != nil {
			fmt.Printf("%#v\n", err)
			return
		}
	} else {
		err = log.NewError("Unknown protocal, supported: mcpack(PROTOCAL_MCPACK), protobuf(PROTOCAL_PROTO)")
		fmt.Printf("%#v\n", err)
	}
	return
}

func autoParse(data []byte, v interface{}) (length int, nh *nshead.NsHead, ch *cshead.CsHead, err error) {
	var (
		body     []byte
		bodyLen  int
		offset   int
		header   int
		protocal int
	)
	// identy header
	header, bodyLen, offset, nh, ch, err = ParseHead(data)
	if err != nil {
		log.Err(UnsupportedHeaderErr)
	} else {
		if header == HEADER_NSHEAD {
			nh, _ = nshead.Unmarshal(data)
			length = int(nh.BodyLen)
			body = data[offset+nshead.NSHEAD_LENGTH : offset+nshead.NSHEAD_LENGTH+length]
		} else if header == HEADER_CSHEAD {
			ch, _ = cshead.Unmarshal(data)
			length = int(ch.BodyLen)
			body = data[offset+cshead.CSHEAD_LENGTH : offset+cshead.CSHEAD_LENGTH+length]
		} else if header == HEADER_INT {
			body = data[4:]
			length = int(bodyLen)
		} else if header == HEADER_INT16 {
			body = data[2:]
			length = int(bodyLen)
		}
	}

	//identy body
	ver, d := mcpack.Identify(body)
	if ver == 1 || ver == 2 {
		if len(body) != len(d) {
			err = log.NewError("Mixed body with mcpack(PROTOCAL_MCPACK) v%d, unknown length=%d", ver, len(body)-len(d))
			fmt.Printf("255 Mixed body with mcpack(PROTOCAL_MCPACK) v%d, unknown length=%d\n", ver, len(body)-len(d))
		}
		body = d
		protocal = PROTOCAL_MCPACK
	} else {
		protocal = PROTOCAL_PROTO
	}

	if protocal == PROTOCAL_MCPACK {
		err = mcpack.Unmarshal(body, v)
		if err != nil {
			fmt.Printf("%#v\n", err)
			return
		}
	} else if protocal == PROTOCAL_PROTO {
		err = proto.Unmarshal(body, v.(proto.Message))
		if err != nil {
			fmt.Printf("%#v\n", err)
			return
		}
	} else {
		err = log.NewError("Unknown protocal, supported: mcpack(PROTOCAL_MCPACK), protobuf(PROTOCAL_PROTO)")
		fmt.Printf("%#v\n", err)
	}
	return
}
func Parser(data []byte) (header, protocal int) {
	f := func(d []byte) error { return output.PrintMcpack(d) }

	var (
		body    []byte
		bodyLen int
		offset  int
		nh      *nshead.NsHead
		ch      *cshead.CsHead
		err     error
	)
	// identy header
	header, bodyLen, offset, nh, ch, err = ParseHead(data)
	if err != nil {
		log.Err(err)
	} else {
		if header == HEADER_NSHEAD {
			bodyLen = int(nh.BodyLen)
			body = data[offset+nshead.NSHEAD_LENGTH : offset+nshead.NSHEAD_LENGTH+bodyLen]
			log.Info("Header: nshead, provider: %s, length: %d", nh.Provider, nh.BodyLen)
		} else if header == HEADER_CSHEAD {
			bodyLen = int(ch.BodyLen)
			body = data[offset+cshead.CSHEAD_LENGTH : offset+cshead.CSHEAD_LENGTH+bodyLen]
			log.Info("Header: cshead, length: %d", ch.BodyLen)
		} else if header == HEADER_INT {
			body = data[4+offset:]
			log.Info("Header: uint32(4 bytes), length: %d", bodyLen)
		} else if header == HEADER_INT16 {
			body = data[2+offset:]
			log.Info("Header: uint16(2 bytes)+(2 bytes), length: %d", bodyLen)
		}
	}

	//identy body
	ver, d := mcpack.Identify(body)
	if ver == 1 || ver == 2 {
		log.Info("Body: mcpack v%d, length: %d", ver, bodyLen)
		if len(body) != len(d) {
			if bodyLen == 0 {
				log.Error("Body: incomplete mcpack(PROTOCAL_MCPACK) v%d data, length=%d", ver, bodyLen)
			} else {
				ulen := len(body) - len(d)
				log.Warning("Body: mixed body with mcpack(PROTOCAL_MCPACK) v%d, length=%d, unknown length=%d(%#v)", ver, bodyLen, ulen, body[:ulen])
			}
		}
		body = d
		f(body)
		protocal = PROTOCAL_MCPACK
	} else {
		log.Error("Unknown protocal, proto maybe. Supported: mcpack(PROTOCAL_MCPACK), protobuf(PROTOCAL_PROTO)")
		protocal = PROTOCAL_PROTO
	}
	return
}
func ParseHead(data []byte) (header int, bodyLen int, offset int, nh *nshead.NsHead, ch *cshead.CsHead, err error) {
	header = HEADER_UNKNOWN
	leng := len(data)
	if leng <= 36 {
		err = log.NewError("Invalid length data")
		return
	}
	if i := bytes.Index(data, nshead.NSHEAD_MN); i != -1 {
		if i > 24 {
			offset = i - 24
			data = data[offset:]
			fmt.Printf("Invalid data, more data[%d bytes] seems be ahead of nshead\n", offset)
		} else if i < 24 {
			//err = log.NewError("Invalid data, maybe duplicated data exist")
			fmt.Printf("Invalid data, nshead is uncompletement[%d bytes]\n", i-24)
		}
		n, ok := nshead.Unmarshal(data)
		if ok {
			header = HEADER_NSHEAD
			nh = n
			bodyLen = int(nh.BodyLen)
			return
		}
		//fmt.Printf("%#v", n)
		if leng != int(n.BodyLen)+nshead.NSHEAD_LENGTH {
			//err = log.NewError("Invalid data, maybe duplicated data exist")
			fmt.Printf("Invalid data, maybe duplicated data exist\n")
		}
	}
	// identy header
	f4headLen := int(binary.LittleEndian.Uint32(data[:4]))
	f2headLen := int(binary.LittleEndian.Uint16(data[:2]))
	headLen := int(binary.LittleEndian.Uint32(data[32:36]))
	log.Info("data length: %d; int32 length: %d; int16 length: %d; nshead/cshead length: %d", leng, f4headLen, f2headLen, headLen)
	if f4headLen+4 == leng {
		header = HEADER_INT
		bodyLen = f4headLen
	} else if f2headLen+2 == leng {
		header = HEADER_INT16
		bodyLen = f2headLen
	} else if f2headLen+4 == leng {
		header = HEADER_INT16
		bodyLen = f2headLen
		offset = 2
		log.Warning("uint16(HEADER_INT16) header, but unknown data existing: %v, offset: %d", data[2:4], offset)
	} else {
		if leng >= 36 {
			c, ok := cshead.Unmarshal(data)
			if ok {
				ch = c
				header = HEADER_CSHEAD
				bodyLen = int(ch.BodyLen)
				return
			}
			if leng != int(c.BodyLen)+cshead.CSHEAD_LENGTH {
				//err = log.NewError("Invalid data, maybe duplicated data exist")
				fmt.Printf("Invalid data, maybe duplicated data exist\n")
			}
		}
		err = UnsupportedHeaderErr
	}
	return
}
