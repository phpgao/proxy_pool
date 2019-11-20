package model

import (
	"fmt"
	"testing"
)

func TestHttpProxy_make(t *testing.T) {
	//var m map[string]string
	m := map[string]string{"Anonymous": "0", "From": "http://www.iphai.com/free/wg", "Ip": "103.194.242.254", "Latency": "2881", "Port": "60025", "Schema": "http", "Score": "50"}
	p, _ := Make(m)
	fmt.Println(p)
}
