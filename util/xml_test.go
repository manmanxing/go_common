package util

import (
	"bytes"
	"fmt"
	"testing"
)

const xmlTest = "<xml><appid>wxddf7e9e1da16c4ad</appid><mch_id>155111111</mch_id><nonce_str>Wuc8iAjANhiNLxjvOCSHeolVgYTY24</nonce_str><sign>303D501EE2F040F5AAB05057CB2B930A</sign><out_refund_no>p_zz_normal_6330721112055402496</out_refund_no><transaction_id></transaction_id><out_trade_no>202010010004535322713576_t</out_trade_no><total_fee>200</total_fee><refund_fee>100</refund_fee><refund_account>REFUND_SOURCE_UNSETTLED_FUNDS</refund_account><bus_name>paymentHub</bus_name><bus_sign>09C150FBEAD0B2A71C13D2C1AF8C602E</bus_sign></xml>"

func TestDecodeXMLToMap(t *testing.T) {
	demo := EncodeXmlToMap(xmlTest)
	fmt.Println(demo)
}

func TestEncodeXMLFromMap(t *testing.T) {
	demo := EncodeXmlToMap(xmlTest)
	w := bytes.NewBufferString("")
	_ = EncodeXMLFromMap(w, demo, "xml")
	fmt.Println(w.String())
}


/**
goos: darwin
goarch: amd64
pkg: github.com/manmanxing/go_common/util
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkXmlStringToMap
BenchmarkXmlStringToMap-4   	   59376	     20123 ns/op
PASS
*/
func BenchmarkXmlStringToMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeXmlToMap(xmlTest)
	}
}

/**
goarch: amd64
pkg: github.com/manmanxing/go_common/util
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkEncodeXMLFromMap
BenchmarkEncodeXMLFromMap-4   	  203290	      6192 ns/op
PASS
*/
func BenchmarkEncodeXMLFromMap(b *testing.B) {
	demo := EncodeXmlToMap(xmlTest)
	for i := 0; i < b.N; i++ {
		w := bytes.NewBufferString("")
		_ = EncodeXMLFromMap(w, demo, "xml")
	}
}