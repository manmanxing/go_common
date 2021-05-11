package handler

import (
	"fmt"
	"time"
)

func Handler() (idle bool, err error) {
	//todo 处理具体的业务逻辑
	fmt.Println("start deal worker")
	time.Sleep(time.Second * 10)
	return false, err
}
