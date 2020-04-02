package net

import "encoding/json"

func GETWithCookie(addr string, cookie string, header ...[2]string) (int, []byte, error) {
	return request("GET", addr, nil, cookie, header...)
}
func GET(addr string, header ...[2]string) (int, []byte, error) {
	return GETWithCookie(addr, "", header...)
}

func GETJSONAndParse(addr string, res interface{}, header ...[2]string) (int, []byte, error) {
	code, resp, err := GET(addr, header...)
	if err != nil {
		return code, nil, err
	}
	return code, resp, json.Unmarshal(resp, &res)
}
