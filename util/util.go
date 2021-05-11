package util

import (
	"encoding/binary"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/manmanxing/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ChangeErr2Grpc(err error) error {
	err = errors.Cause(err)
	if _, ok := status.FromError(err); ok {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return status.Errorf(codes.NotFound, "非法操作")
	}
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return err
}

func GetClientIP(c echo.Context) string {
	ip := c.Get("client_ip")
	if ip == nil {
		addr := c.Request().RemoteAddr
		ip = strings.Split(addr, ":")[0]
		c.Set("client_ip", ip)
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
