package beacon

import (
	"fmt"
	"testing"
)

func TestServiceHost(t *testing.T) {
	ret, err := ServiceHost()
	if err != nil {
		return
	}
	fmt.Println(ret)
}