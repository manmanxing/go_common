package util

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	TimeShortFormat_ = "2006-01-02"
	TimeLongFormat_  = "2006-01-02 15:04:05"
	TimeLongFormat   = "20060102150405"
	TimeShortFormat  = "20060102"
)

var BeijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)

// FormatSequenceNo 格式化 seq 为 12 位的十进制字符串.
func FormatSequenceNo(seq int64) string {
	str := strconv.FormatInt(seq, 10)
	switch n := len(str); {
	case n < 12:
		return "000000000000"[:12-n] + str
	case n == 12:
		return str
	default: // n > 12:
		return str[n-12:]
	}
}

// LowestFourBytesForUserID 获取用户ID的低 4 位.
func LowestFourBytesForUserID(userID int64) string {
	str := strconv.FormatInt(userID, 10)
	switch n := len(str); {
	case n > 4:
		return str[n-4:]
	case n == 4:
		return str
	default: // n < 4:
		return "0000"[:4-n] + str
	}
}

//获取堆栈信息
//要获取全部堆栈信息，可以使用  debug.PrintStack()
func StackInfo() []string {
	var pc [8]uintptr
	sep := "/app/"
	data := make([]string, 0, 8)
	n := runtime.Callers(5, pc[:]) //note
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		if !strings.Contains(file, sep) {
			continue
		}
		ret := strings.Split(file, sep)
		file = ret[1]
		//name := fn.Name()
		data = append(data, fmt.Sprintf("(%s:%d)", file, line))
	}
	return data
}