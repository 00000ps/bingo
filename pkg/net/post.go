package net

import (
	"bingo/pkg/utils"
	"bytes"
	"encoding/json"
	"net/http"
)

// func PostDo(url string, header, forms map[string]string, cookieKey, cookieValue string) (code int, body []byte, err error) {
// 	client := &http.Client{}
// 	if !strings.HasPrefix(url, "http") {
// 		url = "http://" + url
// 	}
// 	var req *http.Request
// 	req, err = http.NewRequest("POST", url, strings.NewReader(""))
// 	//req, err := http.NewRequest("POST", url, strings.NewReader("name=cjb"))
// 	if err != nil {
// 		log.Error("req error: %s", err)
// 		return
// 	}

// 	//add header
// 	if header != nil && len(header) > 0 {
// 		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 		//req.Header.Set("Cookie", cookie)
// 		for k, v := range header {
// 			req.Header.Set(k, v)
// 		}
// 	}
// 	//add cookie
// 	if cookieKey != "" && cookieValue != "" {
// 		expire := time.Now().AddDate(0, 0, 1)
// 		//cookie := &http.Cookie{Name: "BAIDUID", Value: "D31060867A92B7C3C7E637DFC09F09A1:FG=1", Path: "/", Expires: expire, MaxAge: 86400}
// 		cookie := &http.Cookie{Name: cookieKey, Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
// 		req.AddCookie(cookie)
// 	}
// 	//add forms
// 	if forms != nil && len(forms) > 0 {
// 		for k, v := range forms {
// 			req.Form.Add(k, v)
// 		}
// 	}

// 	resp, err := client.Do(req)
// 	code = resp.StatusCode
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return
// 	}

// 	body, err = utils.Read(resp.Body)
// 	if err != nil {
// 		log.Debug("read response error: %s", err)
// 	}
// 	return
// }

func POSTWithCookie(addr string, body []byte, cookie string, header ...[2]string) (int, []byte, error) {
	return request("POST", addr, body, cookie, header...)
}
func POST(addr string, body []byte, header ...[2]string) (int, []byte, error) {
	return POSTWithCookie(addr, body, "", header...)
}

func POSTJSON(addr string, req interface{}, header ...[2]string) (int, []byte, error) {
	d1, e1 := json.Marshal(req)
	if e1 != nil {
		return -1, nil, e1
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(d1))
	if err != nil {
		return resp.StatusCode, nil, err
	}
	defer resp.Body.Close()
	d, e := utils.Read(resp.Body)
	return resp.StatusCode, d, e

	// return POST(addr, d, header...)
}
func POSTJSONAndParse(addr string, req interface{}, res interface{}, header ...[2]string) (int, []byte, error) {
	code, resp, err := POSTJSON(addr, req, header...)
	if err != nil {
		return code, nil, err
	}
	return code, resp, json.Unmarshal(resp, &res)
}
