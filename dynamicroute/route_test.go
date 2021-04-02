package dynamicroute

import (
	"fmt"
	"testing"
)

func TestGetRoute(t *testing.T) {
	path := GetRoute(0)
	fmt.Println(checkRoutePeriods(path, 0))
}

/*
	goos: darwin
	goarch: amd64
	pkg: github.com/manmanxing/go_center_common/dynamicroute
	cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
	BenchmarkGenRoute
	BenchmarkGenRoute-4   	 2218090	       548.8 ns/op
	PASS
*/
func BenchmarkGenRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRoute(0)
	}
}
