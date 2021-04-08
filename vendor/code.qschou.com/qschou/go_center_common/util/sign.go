package util

import (
	"net/url"
	"strings"
)

//微信的md5签名
func GenWxSign(data url.Values, key string) string {
	data.Del("sign")
	for k, _ := range data {
		v := data.Get(k)
		if len(strings.TrimSpace(v)) <= 0 {
			data.Del(k)
		}
	}
	str, _ := url.QueryUnescape(data.Encode())
	str += "&key=" + key
	sign := strings.ToUpper(Md5(str))
	return sign
}

//微众银行的的md5签名
func GenWbSign(data url.Values, key string) string {
	data.Del("sign")
	str, _ := url.QueryUnescape(data.Encode())
	str += "&key=" + key
	sign := strings.ToLower(Md5(str))
	return sign
}
func GenAliMd5Sign(data url.Values, key string) string {
	data.Del("sign")
	data.Del("sign_type")
	str, _ := url.QueryUnescape(data.Encode())
	str += key
	sign := Md5(str)
	return sign
}
