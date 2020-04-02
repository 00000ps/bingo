package encode

import (
	"bingo/pkg/log"
	"bingo/pkg/utils"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/golang/protobuf/proto"
)

const (
	fieldHeader = "Header"
	fieldNeck   = "Neck"
	fieldBody   = "Body"

	tagName     = "json"
	tagProtocal = "protocal"
	tagEncoding = "encoding"
	tagDefault  = "default"
	tagRequired = "required"
	tagType     = "type"

	ProtocalINT16  = "int16"
	ProtocalINT32  = "int32"
	ProtocalNSHEAD = "nshead"
	ProtocalCSHEAD = "cshead"
	ProtocalGET    = "get"
	ProtocalPOST   = "post"
	ProtocalTCP    = "tcp"
	ProtocalUDP    = "udp"
	ProtocalHTTP   = "http"
	ProtocalHTTPS  = "https"

	EncodingJSON       = "json"
	EncodingFormData   = "form-data"
	EncodingURLEncoded = "urlencoded"
	EncodingMCPACK1    = "mcpack1"
	EncodingMCPACK2    = "mcpack2"
	EncodingPROTOBUF   = "proto"
	EncodingSTRUCT     = "struct"
	EncodingBASE64     = "BASE64"
	EncodingUnknown    = "unknown"

	FieldIgnore = "FIELDIGNORED"

	typeFile = "file"

	tagMust     = "must"
	tagOptional = "optional"
	tagFalse    = "false"
)

type Result struct {
	Protocal, Encoding string
	Data               []byte
	http.Header
	EncodingData  []byte // json
	url.Values           // url-encoded
	*bytes.Buffer        // form-data
}

var UnsupportedHeaderErr = fmt.Errorf("Unknown header, supported: int32(HEADER_INT32), int16(HEADER_INT16), http: POST/GET")

func SmartDecoder(data []byte, v interface{}) (err error) {
	// if len(data) < 4 {
	// 	return fmt.Errorf("decoder: invalid data length: %d", len(data))
	// }

	// //默认作为json解析
	// if _, err = utils.JSONPretty(data); err == nil {
	// 	return json.Unmarshal(data, v)
	// }
	// err = nil

	// var body []byte
	// pv := reflect.ValueOf(v)
	// if pv.Kind() != reflect.Ptr || pv.IsNil() {
	// 	return fmt.Errorf("decoder: need Ptr data")
	// }
	// elem := pv.Elem()
	// if elem.Kind() != reflect.Struct {
	// 	return fmt.Errorf("decoder: need Struct data")
	// }

	// for i := 0; i < elem.NumField(); i++ {
	// 	sv := elem.Field(i)
	// 	st := elem.Type().Field(i)
	// 	if st.Tag.Get("json") != "" {
	// 		//json.Unmarshal(data, v)
	// 		return errors.New("json Parsed Unmatched")
	// 	}
	// 	if st.Type.Kind() == reflect.Struct {
	// 		protocal := st.Type.Name()
	// 		if protocal == "" {
	// 			protocal = st.Tag.Get(tagProtocal)
	// 		}
	// 		//fmt.Printf("protocal:%s; val=%s\n\n", st.Tag.Get(tagProtocal), st.Tag.Get("val"))
	// 		protocal = strings.ToLower(protocal)
	// 		//fmt.Printf("%#v--name: %#v--protocal: %#v--type: %#v\n", st, st.Name, protocal, st.Type.Name())
	// 		if st.Name == fieldHeader {
	// 			if protocal == "nshead" {
	// 				n, _ := nshead.Unmarshal(data)
	// 				sv.Set(reflect.ValueOf(*n))
	// 				body = data[nshead.NSHEAD_LENGTH:]
	// 			} else if protocal == "cshead" {
	// 				c, _ := cshead.Unmarshal(data)
	// 				sv.Set(reflect.ValueOf(*c))
	// 				body = data[cshead.CSHEAD_LENGTH:]
	// 			} else {
	// 				log.Error("protocal: %s; %s", protocal, UnsupportedHeaderErr)
	// 			}
	// 			//fmt.Printf("%#v\n", sv)
	// 		} else if st.Name == fieldBody {
	// 			if strings.Contains(protocal, "mcpack") {
	// 				err = mcpack.Unmarshal(body, sv)
	// 				if err != nil {
	// 					//fmt.Printf("%#v\n", sv)
	// 					fmt.Printf("%#v\n", err)
	// 					return
	// 				}
	// 				//fmt.Printf("%#v\n", sv)
	// 			} else if strings.Contains(protocal, "proto") {
	// 				err = proto.Unmarshal(body, sv.Interface().(proto.Message))
	// 				if err != nil {
	// 					//fmt.Printf("%#v\n", sv)
	// 					fmt.Printf("%#v\n", err)
	// 					return
	// 				}
	// 				//fmt.Printf("%#v\n", sv)
	// 			}
	// 		} else if st.Name == fieldNeck {
	// 			if strings.Contains(protocal, "struct") {
	// 				structs.Unmarshal(body, sv)
	// 				//fmt.Printf("%#v\n", sv)
	// 			}
	// 		}
	// 	} else if st.Type.Kind() == reflect.Uint32 || st.Type.Kind() == reflect.Int32 || st.Type.Kind() == reflect.Int {
	// 		if st.Name == fieldHeader {
	// 			n := binary.LittleEndian.Uint32(data[0:4])
	// 			if st.Type.Kind() == reflect.Uint32 {
	// 				sv.Set(reflect.ValueOf(n))
	// 			} else if st.Type.Kind() == reflect.Int32 {
	// 				sv.Set(reflect.ValueOf(int32(n)))
	// 			} else if st.Type.Kind() == reflect.Int {
	// 				sv.Set(reflect.ValueOf(int(n)))
	// 			}
	// 			body = data[4:]
	// 			//fmt.Printf("%#v\n", sv)
	// 		} else if st.Name == fieldNeck {
	// 			n := binary.LittleEndian.Uint32(body[0:4])
	// 			if st.Type.Kind() == reflect.Uint16 {
	// 				sv.Set(reflect.ValueOf(n))
	// 			} else if st.Type.Kind() == reflect.Int32 {
	// 				sv.Set(reflect.ValueOf(int32(n)))
	// 			} else if st.Type.Kind() == reflect.Int {
	// 				sv.Set(reflect.ValueOf(int(n)))
	// 			}
	// 			body = body[4:]
	// 		} else {
	// 			log.Error("st.Name: %s; %s", st.Name, UnsupportedHeaderErr)
	// 		}
	// 	} else if st.Type.Kind() == reflect.Uint16 || st.Type.Kind() == reflect.Int16 {
	// 		if st.Name == fieldHeader {
	// 			n := binary.LittleEndian.Uint16(data[0:2])
	// 			if st.Type.Kind() == reflect.Uint16 {
	// 				sv.Set(reflect.ValueOf(n))
	// 			} else if st.Type.Kind() == reflect.Int16 {
	// 				sv.Set(reflect.ValueOf(int16(n)))
	// 			}
	// 			body = data[2:]
	// 			//fmt.Printf("%#v\n", sv)
	// 		} else if st.Name == fieldNeck {
	// 			n := binary.LittleEndian.Uint16(body[0:2])
	// 			if st.Type.Kind() == reflect.Uint16 {
	// 				sv.Set(reflect.ValueOf(n))
	// 			} else if st.Type.Kind() == reflect.Int16 {
	// 				sv.Set(reflect.ValueOf(int16(n)))
	// 			}
	// 			body = body[2:]
	// 		} else {
	// 			log.Error("st.Name: %s; %s", st.Name, UnsupportedHeaderErr)
	// 		}
	// 	} else if st.Type.Kind() == reflect.Ptr {
	// 		elem := st.Type.Elem()
	// 		if elem.Kind() != reflect.Struct {
	// 			err = errors.New("parser need Struct data(not pointer)")
	// 			return
	// 		}
	// 	}
	// }

	// if len(data) > 0 {
	// 	a := mcpack.Decode(data)
	// 	// TODO: complete SmartEncoder
	// 	out, _ := SmartEncoder(v)
	// 	b := mcpack.Decode(out.Data)
	// 	if !utils.Assert(a, b) {
	// 		ts := "log/" + utils.GetTimestamp()
	// 		rfile := ts + "_raw"
	// 		nfile := ts + "_new"
	// 		log.Warning("encode warning! interface mismatched, please check: [vimdiff %s %s]", rfile, nfile)
	// 		utils.WriteNewFile(rfile, []byte(a))
	// 		utils.WriteNewFile(nfile, []byte(b))
	// 		//fmt.Printf(color.Green("interface RAW: ========\n%s\n", a))
	// 		//fmt.Printf(color.Red("interface NEW: ========\n%s\n", b))
	// 	}
	// }

	return
}
func SmartEncoder(v interface{}) (out *Result, msg utils.Msg) {
	// out = new(Result)
	// //msg = utils.Msg{}
	// var (
	// 	header           string
	// 	nh               *nshead.NsHead
	// 	ch               *cshead.CsHead
	// 	head, neck, body []byte
	// 	err              error
	// )
	// pv := reflect.ValueOf(v)
	// if pv.Kind() != reflect.Ptr || pv.IsNil() {
	// 	msg.SetError("encoder need Ptr data")
	// 	return
	// }
	// elem := pv.Elem()
	// if elem.Kind() != reflect.Struct {
	// 	msg.SetError("encoder need Struct data")
	// 	return
	// }

	// for i := 0; i < elem.NumField(); i++ {
	// 	sv := elem.Field(i)
	// 	st := elem.Type().Field(i)
	// 	//value := st.Tag.Get("value")
	// 	switch st.Type.Kind() {
	// 	case reflect.Struct:
	// 		if out.Protocal == "" {
	// 			//protocal = st.Type.Name()
	// 			out.Protocal = strings.ToLower(st.Tag.Get(tagProtocal))
	// 		}
	// 		//fmt.Printf("%#v--name: %#v--protocal: %#v--type: %#v\n", st.Type.Name(), st.Name, protocal, st.Type.Name())
	// 		switch st.Name {
	// 		case fieldHeader:
	// 			header = out.Protocal
	// 			switch out.Protocal {
	// 			case ProtocalNSHEAD:
	// 				n := sv.Interface().(nshead.NsHead)
	// 				nh = &n
	// 			case ProtocalCSHEAD:
	// 				c := sv.Interface().(cshead.CsHead)
	// 				ch = &c
	// 			case ProtocalGET, ProtocalPOST:
	// 				out.Header = make(http.Header)
	// 				for i := 0; i < sv.NumField(); i++ {
	// 					required := sv.Type().Field(i).Tag.Get(tagRequired)
	// 					key := sv.Type().Field(i).Tag.Get(tagName)
	// 					val := sv.Field(i).String()
	// 					if val == "" {
	// 						val = sv.Type().Field(i).Tag.Get(tagDefault)
	// 					}
	// 					toSet := false
	// 					if required == tagMust {
	// 						toSet = true
	// 					} else if required == tagOptional {
	// 						if val != "" {
	// 							toSet = true
	// 						}
	// 					} else if required == tagFalse {
	// 						if sv.Field(i).String() != "" {
	// 							msg.SetWarning("abandoned field. %s(%s) has been set as %s: %s", key, tagRequired, tagFalse, sv.Field(i).String())
	// 						}
	// 					} else {
	// 						toSet = true
	// 					}

	// 					if toSet {
	// 						out.Header.Add(key, val)
	// 					}
	// 				}

	// 			default:
	// 				log.Debug("%#v", msg.SetErr(UnsupportedHeaderErr))
	// 			}
	// 		case fieldBody:
	// 			switch out.Protocal {
	// 			case EncodingMCPACK1:
	// 				body, err = mcpack.MarshalV1(sv)
	// 				if err != nil {
	// 					msg.SetErr(err)
	// 					return
	// 				}
	// 			case EncodingMCPACK2:
	// 				body, err = mcpack.MarshalV2(sv)
	// 				if err != nil {
	// 					msg.SetErr(err)
	// 					return
	// 				}
	// 			case EncodingPROTOBUF:
	// 				body, err = proto.Marshal(sv.Interface().(proto.Message))
	// 				if err != nil {
	// 					msg.SetErr(err)
	// 					return
	// 				}

	// 			case ProtocalGET, ProtocalPOST:
	// 				out.Encoding = st.Tag.Get(tagEncoding)
	// 				//log.Notice("%#v----%#v", out, sv.Interface())
	// 				switch out.Encoding {
	// 				case EncodingJSON:
	// 					out.EncodingData, err = json.Marshal(sv.Interface())
	// 					if err != nil {
	// 						msg.SetErr(err)
	// 					}
	// 					out.Header = make(http.Header)
	// 					out.Header.Set("Content-Type", "application/json")
	// 					return
	// 				case EncodingFormData:
	// 					b, w, err := postFormData(sv)
	// 					if err != nil {
	// 						msg.SetErr(err)
	// 					}
	// 					out.Header = make(http.Header)
	// 					out.Header.Set("Content-Type", w.FormDataContentType())
	// 					out.Buffer = b
	// 					out.Data = b.Bytes()
	// 					return

	// 				default:
	// 					out.Values = url.Values{}
	// 					if war := postEncode(out.Values, sv); war != "" {
	// 						msg.SetWarning(war)
	// 					}
	// 					return
	// 				}
	// 			default:
	// 				//log.Notice("eeeeeeeeee: %s", protocal)
	// 			}
	// 		case fieldNeck:
	// 			if out.Protocal == EncodingSTRUCT {

	// 			}
	// 		}

	// 	case reflect.Slice:
	// 		if out.Protocal == "" {
	// 			out.Protocal = strings.ToLower(st.Tag.Get(tagProtocal))
	// 		}
	// 		switch st.Name {
	// 		case fieldHeader:
	// 		case fieldBody:
	// 			switch out.Protocal {
	// 			case EncodingMCPACK1:
	// 			case EncodingMCPACK2:
	// 			case EncodingPROTOBUF:
	// 			case ProtocalGET, ProtocalPOST:
	// 				out.Encoding = st.Tag.Get(tagEncoding)
	// 				switch out.Encoding {
	// 				case EncodingJSON:
	// 					//for i := 0; i < sv.Len(); i++ {
	// 					//	sv.Index(i)
	// 					//}
	// 					out.EncodingData, err = json.Marshal(sv.Interface())
	// 					if err != nil {
	// 						msg.SetErr(err)
	// 					}
	// 					return
	// 				default:
	// 				}
	// 			default:
	// 				//log.Notice("eeeeeeeeee: %s", protocal)
	// 			}

	// 		}
	// 	case reflect.Uint32, reflect.Int32:
	// 		if st.Name == fieldHeader {
	// 			header = ProtocalINT32
	// 		} else if st.Name == fieldNeck {
	// 			b := bytes.NewBuffer(make([]byte, 0, 4))
	// 			binary.Write(b, binary.LittleEndian, sv.Interface())
	// 			neck = b.Bytes()
	// 		} else {
	// 			log.Debug("%#v", msg.SetErr(UnsupportedHeaderErr))
	// 		}
	// 	case reflect.Uint16, reflect.Int16:
	// 		if st.Name == fieldHeader {
	// 			header = ProtocalINT16
	// 		} else if st.Name == fieldNeck {
	// 			b := bytes.NewBuffer(make([]byte, 0, 2))
	// 			binary.Write(b, binary.LittleEndian, sv.Interface())
	// 			neck = b.Bytes()
	// 		} else {
	// 			log.Debug("%#v", msg.SetErr(UnsupportedHeaderErr))
	// 		}
	// 	case reflect.Ptr:
	// 		elem := st.Type.Elem()
	// 		if elem.Kind() != reflect.Struct {
	// 			msg.SetError("parser need Struct data(not pointer)")
	// 			return
	// 		}
	// 	case reflect.String:
	// 		if out.Protocal == "" {
	// 			out.Protocal = strings.ToLower(st.Tag.Get(tagProtocal))
	// 		}
	// 		if st.Name == fieldHeader {
	// 			header = out.Protocal
	// 		}
	// 	case reflect.Map:
	// 		out.Values = url.Values{}
	// 		for i := 0; i < sv.NumField(); i++ {
	// 			if sv.Field(i).String() != "" {
	// 				out.Values.Add(sv.Type().Field(i).Name, sv.Field(i).String())
	// 			}
	// 		}
	// 		return
	// 	}
	// }

	// if out.Protocal != ProtocalGET && out.Protocal != ProtocalPOST {
	// 	if len(body) < 4 {
	// 		msg.SetError("failed to encode body")
	// 		return
	// 	}

	// 	switch header {
	// 	case ProtocalNSHEAD:
	// 		if nh == nil {
	// 			msg.SetError("failed to encode nshead")
	// 			return
	// 		}
	// 		if len(body) < nshead.NSHEAD_LENGTH {
	// 			msg.SetError("failed to encode nshead body")
	// 			return
	// 		}
	// 		head = nh.Encode(body)
	// 		if len(head) == 0 {
	// 			msg.SetError("failed to encode nshead body")
	// 			return
	// 		}
	// 	case ProtocalCSHEAD:
	// 		if ch == nil {
	// 			msg.SetError("failed to encode cshead")
	// 			return
	// 		}
	// 		if len(body) < cshead.CSHEAD_LENGTH {
	// 			msg.SetError("failed to encode cshead body")
	// 			return
	// 		}
	// 		head = ch.Encode(body)
	// 		if len(head) == 0 {
	// 			msg.SetError("failed to encode cshead body")
	// 			return
	// 		}
	// 	case ProtocalINT32:
	// 		b := bytes.NewBuffer(make([]byte, 0, 4))
	// 		binary.Write(b, binary.LittleEndian, uint32(len(body)))
	// 		head = b.Bytes()
	// 	case ProtocalINT16:
	// 		b := bytes.NewBuffer(make([]byte, 0, 2))
	// 		binary.Write(b, binary.LittleEndian, uint16(len(body)))
	// 		head = b.Bytes()
	// 	}

	// 	out.Data = append(head, neck...)
	// 	out.Data = append(out.Data, body...)
	// }
	return
}

func reflectField(v reflect.Value) (val string, zero bool) {
	switch v.Kind() {
	case reflect.Float32:
		val = strconv.FormatFloat(v.Float(), 'f', 10, 32)
	case reflect.Float64:
		val = strconv.FormatFloat(v.Float(), 'f', 10, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = strconv.FormatUint(v.Uint(), 10)
	case reflect.String:
		val = v.String()
	case reflect.Slice, reflect.Array:
		val = string(v.Bytes())
	case reflect.Ptr:
		// if v.Elem() != reflect.Zero() {
		val, zero = reflectField(v.Elem())
		// }
		// pv := reflect.ValueOf(svv)
		// pt := reflect.TypeOf(svv)
		// log.PrintRaw(v.Elem().Type().String())
		// log.PrintRaw(v.Elem())
	case reflect.Invalid:
		zero = true
	default:
		log.Debug("unexpected type: %d=%v", v.Kind(), v.Kind())
	}
	return
}
func postEncode(uv url.Values, sv reflect.Value) (war string) {
	for i := 0; i < sv.NumField(); i++ {
		//required := sv.Type().Field(i).Tag.Get(tagRequired)
		svv := sv.Field(i)
		val, zero := reflectField(svv)

		key := sv.Type().Field(i).Tag.Get(tagName)
		oe := strings.HasSuffix(key, ",omitempty")
		if oe {
			key = strings.TrimSuffix(key, ",omitempty")
		}

		if (val == FieldIgnore) || (oe && zero) {
			log.Debug("field %s=%s filted", key, val)
			continue
		}
		uv.Add(key, val)
	}
	return war
}
func postFormData(sv reflect.Value) (b *bytes.Buffer, w *multipart.Writer, err error) {
	//log.PrintRaw(sv)
	b = new(bytes.Buffer)
	w = multipart.NewWriter(b)
	defer w.Close()

	for i := 0; i < sv.NumField(); i++ {
		//required := sv.Type().Field(i).Tag.Get(tagRequired)
		key := sv.Type().Field(i).Tag.Get(tagName)
		tType := sv.Type().Field(i).Tag.Get(tagType)

		var val string
		// var fw io.Writer

		switch sv.Field(i).Kind() {
		case reflect.String:
			val = sv.Field(i).String()
			if val != FieldIgnore {
				switch tType {
				case typeFile:
					fw, err0 := w.CreateFormFile(key, val)
					if err0 != nil {
						err = fmt.Errorf("failed to encode form-data body: %s", err0)
						return
					}
					d, _ := utils.ReadFile(val)
					fw.Write(d)
				default:
					fw, err0 := w.CreateFormField(key)
					if err0 != nil {
						err = fmt.Errorf("failed to encode form-data body: %s", err0)
						return
					}
					if _, err = fw.Write([]byte(val)); err != nil {
						err = fmt.Errorf("failed to encode form-data body: %s", err)
						return
					}
				}
			}
		default:
			log.Debug("unexpected type: %v", sv.Field(i).Kind())
		}
	}
	return
}
