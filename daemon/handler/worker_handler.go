package handler

import (
	"fmt"
	"time"
)

//idle:为 true 表示本次没有需要处理的数据
func Handler() (idle bool, err error) {
	//todo 处理具体的业务逻辑
	fmt.Println("start deal worker")
	time.Sleep(time.Second * 5)
	return false, err
}
