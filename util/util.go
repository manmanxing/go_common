package util

import (
	"strconv"
)

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
