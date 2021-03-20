package http

import (
	"net"
	"net/http"
	"time"
)

var Default = NewHttpClient(time.Second * 5)
var DefaultWithConcurrent = NewConcurrentHttpClient(time.Second*5, 100, 50)

/*
var DefaultTransport RoundTripper = &Transport{
	Proxy: ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
*/

func NewHttpClient(timeOut time.Duration) *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeOut,
	}
}

func NewConcurrentHttpClient(timeOut time.Duration, MaxIdleConn, maxIdleConnPerHost int) *http.Client {
	demoTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 30,
			KeepAlive: time.Second * 30,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          MaxIdleConn,
		MaxIdleConnsPerHost:   maxIdleConnPerHost, //note: default is 2
		IdleConnTimeout:       time.Second * 90,
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: time.Second * 1,
	}
	return &http.Client{
		Transport: demoTransport,
		Timeout:   timeOut,
	}
}