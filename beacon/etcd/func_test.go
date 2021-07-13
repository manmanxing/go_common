package etcd

import (
	"fmt"
	"testing"
)

func TestFunc(t *testing.T) {
	_,err := NewClient()
	if err != nil {
		fmt.Println("err:",err)
		return
	}
}