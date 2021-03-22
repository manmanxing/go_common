package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/manmanxing/errors"
	"github.com/manmanxing/go_center_common/util"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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

type HttpClient struct {
	http.Client
}

func (h *HttpClient) HttpGet(url string, data url.Values) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	if len(data) > 0 {
		url += "?" + data.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = errors.Wrap(err)
		return
	}
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()
	r, err := h.Do(req)
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	return h.getStrFromResp(r)
}

func (h *HttpClient) HttpPost(url string, bodyType string, body string) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	if len(strings.TrimSpace(bodyType)) <= 0 {
		return "", errors.New("bodyType is empty")
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()

	req.Header.Set("Content-Type", bodyType)
	r, err := h.Do(req)
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	return h.getStrFromResp(r)
}

func (h *HttpClient) HttpPut(url string, bodyType string, body string) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	if len(strings.TrimSpace(bodyType)) <= 0 {
		return "", errors.New("bodyType is empty")
	}

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()

	req.Header.Set("Content-Type", bodyType)
	r, err := h.Do(req)
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	return h.getStrFromResp(r)
}

func (h *HttpClient) HttpPostJson(url string, data interface{}) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	body := ""
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	switch dataValue.Kind() {
	case reflect.String:
		body = data.(string)
	case reflect.Map, reflect.Struct:
		tmp, err := json.Marshal(data)
		if err != nil {
			err = errors.Wrap(err)
			return "", err
		}
		body = string(tmp)
	default:
		return "", errors.Errorf("unsupported data type: %s", dataValue.Kind().String())
	}

	return h.HttpPost(url, "application/json;charset=utf-8", body)
}

func (h *HttpClient) HttpPostForm(url string, data url.Values) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	body := data.Encode()
	return h.HttpPost(url, "application/x-www-form-urlencoded", body)
}

func (h *HttpClient) HttpPostXml(url string, data url.Values) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	dt := make(map[string]string)
	for k := range data {
		dt[k] = data.Get(k)
	}
	w := bytes.NewBufferString("")
	err = util.EncodeXMLFromMap(w, dt, "xml")
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	return h.HttpPost(url, "application/xml", w.String())
}

func (h *HttpClient) HttpPostWithHeader(url string, headers map[string]string, body string) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	r, err := h.Do(req)
	if err != nil {
		err = errors.Wrap(err)
		return "", err
	}
	return h.getStrFromResp(r)
}

func (h *HttpClient) HttpPutJson(url string, data interface{}) (resp string, err error) {
	if len(strings.TrimSpace(url)) <= 0 {
		return "", errors.New("url is empty")
	}
	body := ""
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	switch dataValue.Kind() {
	case reflect.String:
		body = data.(string)
	case reflect.Map, reflect.Struct:
		tmp, err := json.Marshal(data)
		if err != nil {
			err = errors.Wrap(err)
			return "", err
		}
		body = string(tmp)
	default:
		return "", errors.Errorf("unsupported data type: %s", dataValue.Kind().String())
	}

	return h.HttpPut(url, "application/json;charset=utf-8", body)
}

func (h *HttpClient) getStrFromResp(resp *http.Response) (str string, err error) {
	if resp.Body == nil {
		return "", nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	str = strings.Trim(string(data), "\r\n")
	if resp.StatusCode != http.StatusOK {
		return str, errors.Errorf("statusCode=%v;resp=%s", resp.StatusCode, str)
	}
	return str, nil
}

//增加 http 耗时日志打印
func (h *HttpClient) log(urlPath, body, resp string, dur time.Duration, err error) {
	urlPath, _ = url.QueryUnescape(urlPath)
	resp = strings.Replace(resp, "\r", "", -1)
	resp = strings.Replace(resp, "\n", " ", -1)
	info := fmt.Sprintf("url %s,body %s,resp %s,err %v,time(ms) %d \n", urlPath, body, resp, err, int64(dur/time.Millisecond))
	if err != nil {
		log.Error(info)
	} else {
		log.Info(info)
	}
	if dur > time.Millisecond*500 {
		log.Warn(info)
	}
}