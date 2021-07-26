package beacon

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/manmanxing/errors"
)

const BEACON_HOST = "BEACON_HOST"

var __httpClient = &http.Client{
	Timeout: time.Second * 5,
}

func getBeaconHost() (host string, err error) {
	defer func() {
		fmt.Println("beaconHost", host, "error", err)
	}()
	v, ok := os.LookupEnv(BEACON_HOST)
	if !ok {
		err = errors.Wrap(errors.New("environment variable not found: BEACON_HOST"))
		return "", err
	}
	return strings.TrimSpace(v), nil
}

//返回服务对外提供服务的 grpc  host
func ServiceHost() (ret string, err error) {
	defer func() {
		fmt.Println("ServiceHost", ret, "error", err)
	}()
	host, err := getBeaconHost()
	if err != nil {
		return "", err
	}
	resp, err := __httpClient.Get(fmt.Sprintf("http://%s/service_host", host))
	if err != nil {
		err = errors.Wrap(err, "get service host err")
		return "", err
	}
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			err = errors.Wrap(e, "http resp body close err")
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(fmt.Errorf("get service host failed, returned http status: %s", resp.Status))
		return "", err
	}
	//body: a0e5e536e0eca350efc7c57908b0ea69127.0.0.1
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "io read err")
		return "", err
	}
	body = bytes.TrimSpace(body)
	bodyStr := string(body)
	fmt.Println(bodyStr)
	sumLen := hex.EncodedLen(md5.Size)
	if len(body) <= sumLen {
		err = errors.Wrap(fmt.Errorf("get service host failed, too short, returned http body: %s", body))
		return "", err
	}
	//haveSum: a0e5e536e0eca350efc7c57908b0ea69
	haveSum := body[:sumLen]
	//ip: 127.0.0.1
	ip := body[sumLen:]

	wantSum := make([]byte, sumLen)
	sum := md5.Sum(ip)
	hex.Encode(wantSum, sum[:])
	if !bytes.Equal(haveSum, wantSum) {
		err = errors.Wrap(fmt.Errorf("get service host failed, hashsum mismatch, returned http body: %s", body))
		return "", err
	}

	return string(ip), nil
}
