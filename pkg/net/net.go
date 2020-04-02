package net

import (
	"bingo/pkg/utils"
	"bytes"
	"net/http"
	"strings"
	"time"
)

func request(method string, addr string, body []byte, cookie string, header ...[2]string) (int, []byte, error) {
	var req *http.Request
	if body != nil && len(body) > 0 {
		req, _ = http.NewRequest(method, addr, bytes.NewReader(body))
	} else {
		req, _ = http.NewRequest(method, addr, nil)
	}
	// c:=req.
	if cookie != "" {
		raw := cookie
		arr := strings.Split(raw, "; ")
		for _, cs := range arr {
			// kv := strings.Split(cs, "=")
			i := strings.Index(cs, "=")
			k := cs[:i]
			v := cs[i+1:]
			c := &http.Cookie{
				Name:    k,
				Value:   v,
				Path:    "/",
				Expires: time.Now().AddDate(0, 1, 0),
				MaxAge:  86400}
			// log.Notice("k=%s; v=%s", k, v)
			req.AddCookie(c)
		}
	}
	if header != nil {
		for _, h := range header {
			req.Header.Set(h[0], h[1])
		}
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	// log.Debug(req.URL.String())
	// log.Notice("%#v", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// log.Debug("net request err1: %s", err)
		if resp == nil {
			return -1, nil, err
		}
		return resp.StatusCode, nil, err
	}
	result, err := utils.Read(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		// log.Debug("net request err2: %s", err)
		return resp.StatusCode, nil, err
	}
	// log.Debug("%s", result)
	return resp.StatusCode, result, nil
}
