package util

import (
	"strings"
)

// 等待响应头超时
func IsNetHttpTimeoutAwaitingHeaders(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "net/http: request canceled (Client.Timeout exceeded while awaiting headers)") { // net/http: request canceled (Client.Timeout exceeded while awaiting headers)
		return true
	}
	return false
}

// 读取响应体超时
func IsNetHttpTimeoutReadingBody(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "net/http: request canceled (Client.Timeout exceeded while reading body)") { // net/http: request canceled (Client.Timeout exceeded while reading body)
		return true
	}
	return false
}

// 建立连接超时
func IsNetHttpTimeoutWaitingForConnection(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)") { // net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
		return true
	}
	return false
}
