package util

import (
	"crypto/md5"
	"encoding/hex"
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

func Md5(str string) string {
	tmp := md5.Sum([]byte(str))
	return hex.EncodeToString(tmp[:])
}
