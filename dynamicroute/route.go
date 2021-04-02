package dynamicroute

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	expireTimeRange    = 60 * 10            //动态路由有效时间，单位 秒
	dynamicRouteLength = 4                  //动态路由长度
	seed               = "azxcvbnmsdfghjkl" //随机数种子
)

//now : 当前时间戳
func GetRoute(now int64) string {
	if now <= 0 {
		now = time.Now().Unix()
	}
	c := now / expireTimeRange
	return sum(c)
}

func sum(c int64) string {
	s := md5.Sum([]byte(fmt.Sprintf("%d%s", c, seed)))
	var encrypt []byte
	position := c % md5.Size
	for i := 0; i < dynamicRouteLength; i++ {
		b := s[position%md5.Size]
		encrypt = append(encrypt, b)
		position = position + 1
	}
	return hex.EncodeToString(encrypt)
}

//route:动态路由地址
//period:时间段，即  expireTimeRange * period
func checkRoutePeriods(route string, period int) (result bool) {
	if len(route) != 2*dynamicRouteLength {
		return
	}
	now := time.Now().Unix()
	count := now / expireTimeRange
	//优先考虑当前时间的一个周期，然后再是上几个时间周期
	for i := 0; i < period+1; i++ {
		next := count - int64(i)
		s := sum(next)
		if s == route {
			result = true
			return
		}
	}
	//最后再考虑下一个时间周期
	s2 := sum(count + 1)
	result = s2 == route
	return
}
