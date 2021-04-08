package util

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"crypto/md5"

	"github.com/chanxuehong/uuid"
	"github.com/labstack/echo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 获取 session 中的 user_id  这个方法移植到 go_center_common/api/middleware/auth/get_session_info.go 中 GetSessionUserInfo
//func GetSessionInfo(c echo.Context) *auth.SessionUserInfo {
//	return c.Get(auth.SessionUserInfoContextKey).(*auth.SessionUserInfo)
//}
func GetClientIP(c echo.Context) string {
	ip := c.Get("clientIP")
	if ip == nil {
		addr := c.Request().RemoteAddr
		ip = strings.Split(addr, ":")[0]
		c.Set("clientIP", ip)
	}
	return ip.(string)
}

func GenerateSpanID(addr string) string {
	strAddr := strings.Split(addr, ":")
	ip := strAddr[0]
	ipLong, _ := Ip2Long(ip)
	times := uint64(time.Now().UnixNano())
	spanId := ((times ^ uint64(ipLong)) << 32) | uint64(rand.Int31())
	return strconv.FormatUint(spanId, 16)
}

func Ip2Long(ip string) (uint32, error) {
	ipAddr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(ipAddr.IP.To4()), nil
}

func RandNumStr(len uint8) string {
	r := NewRandom()
	r.SetCharset(Numeric)
	return r.String(len)
}

// 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func GUID() string {
	return uuid.NewV1().String()
	//由时间戳 8字节, 本地ip 四字节, 自增序列 8字节,会轮询,保证 唯一
	baseNum := uint64(0)
	ch := make(chan uint64, 100)
	go func(ch chan uint64) {
		for {
			baseNum++
			if baseNum >= math.MaxUint64 {
				baseNum = 0
			}
			ch <- baseNum
		}
	}(ch)
	num := <-ch
	byNum := make([]byte, 8)
	binary.BigEndian.PutUint64(byNum, num)
	byTime := make([]byte, 8)
	binary.BigEndian.PutUint64(byTime, uint64(time.Now().UnixNano()))
	byIP := []byte(GetLocalIP().To4())
	data := make([]byte, 16, 20)
	for i := 0; i < 8; i++ {
		data[i*2] = byNum[i]
		data[i*2+1] = byTime[i]
	}
	tmp := append(data[0:8], byIP...)
	tmp = append(tmp, data[8:16]...)
	return hex.EncodeToString(tmp)

}
func GetLocalIP() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IPv4zero
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				return ipnet.IP
			}
		}
	}
	return net.IPv4zero
}
func GetLocalPublicIP() net.IP {
	return GetLocalWanIP()
}
func GetLocalWanIP() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IPv4zero
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil && publicIP(ip) {
				return ipnet.IP
			}
		}
	}
	return net.IPv4zero
}

//pre 表示局域ip的第一个字节，10 或者172 或者192
func GetLocalLanIP(pre byte) net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IPv4zero
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil && ip[0] == pre {
				return ipnet.IP
			}
		}
	}
	return net.IPv4zero
}

//判断是否是公网ip
func publicIP(ip net.IP) bool {
	return ip[0] != 10 && ip[0] != 172 && ip[0] != 192
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		if strings.Contains(name, ",") {
			tmp := strings.Split(name, ",")
			name = tmp[0]
		}
		if len(name) <= 0 {
			name = t.Field(i).Name
		}
		data[name] = v.Field(i).Interface()
	}
	return data
}
func Md5(str string) string {
	tmp := md5.Sum([]byte(str))
	return hex.EncodeToString(tmp[:])
}

func ShortFilePath(path string) string {
	ret := strings.Split(path, "qschou")
	if len(ret) >= 2 {
		return ret[1]
	}
	return ret[0]
}

func IsLowerEighteen(IdNumber string) bool {

	if len(IdNumber) < 18 {
		return false
	}
	birth, err := time.Parse("20060102", IdNumber[6:14])
	if err != nil {
		return false
	}
	diff := (time.Now().Unix() - birth.Unix()) / 86400 / 365 //周岁算法.

	if diff < 18 && diff >= 0 {
		return true
	}
	return false

}

func GetAge(IdNumber string) int64 {

	if len(IdNumber) < 18 {
		return 0
	}
	birth, err := time.Parse("20060102", IdNumber[6:14])
	if err != nil {
		return 0
	}
	diff := (time.Now().Unix() - birth.Unix()) / 86400 / 365 //周岁算法.

	if diff <= 10 {
		return 0
	} else if diff > 10 && diff <= 60 {
		return 1
	} else if diff > 60 && diff <= 70 {
		return 2
	} else {
		return 3
	}

}

func GetFee(Coverage, Age int64) int64 {

	switch Coverage {

	case 100000:
		Coverage = 0
	case 200000:
		Coverage = 1
	case 300000:
		Coverage = 2
	default:
		Coverage = 3
	}

	var Fee = [][]int64{
		{3100, 6200, 9300},
		{1500, 3000, 4500},
		{3100, 6200, 9300},
	}

	if Coverage > 2 || Coverage < 0 || Age < 0 || Age > 2 {
		return 0
	}
	return Fee[Age][Coverage]
}
