package model

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type HttpProxy struct {
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Schema    string `json:"schema"`
	Score     int    `json:"score"`
	Latency   int    `json:"latency"`
	From      string `json:"-"`
	Anonymous int    `json:"anonymous"`
}

func Make(m map[string]string) (newProxy HttpProxy, err error) {
	rVal := reflect.ValueOf(&newProxy).Elem()
	rType := reflect.TypeOf(newProxy)
	fieldCount := rType.NumField()

	for i := 0; i < fieldCount; i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		if v, ok := m[t.Name]; ok {
			ddd := reflect.TypeOf(v)
			if ddd != t.Type {
				v, _ := strconv.Atoi(v)
				f.Set(reflect.ValueOf(v))
			} else {
				f.Set(reflect.ValueOf(v))
			}
		} else {
			return newProxy, errors.New(t.Name + " not found")
		}
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	return
}

func (p *HttpProxy) GetKey() string {
	hash := md5.New()
	_, err := io.WriteString(hash, p.GetProxyUrl())
	if err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (p *HttpProxy) GetProxyUrl() string {
	return fmt.Sprintf("%s:%s", p.Ip, p.Port)
}
func (p *HttpProxy) GetProxyHash() map[string]interface{} {
	return structs.Map(p)
}

func (p *HttpProxy) GetIp() string {
	return p.Ip
}

func (p *HttpProxy) GetPort() string {
	return p.Port
}

func (p *HttpProxy) SimpleTcpTest() bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", p.Ip, p.Port), 3*time.Second)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	if err != nil {
		return false
	} else {
		return true
	}
}

func (p *HttpProxy) GetHttpTransport() (t *http.Transport, err error) {
	proxyUrl := &url.URL{Host: p.GetProxyUrl(), Scheme: "http"}

	t = &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	return
}

func (p *HttpProxy) TestProxy(https bool) (err error) {

	startAt := time.Now()
	timeout := 6 * time.Second

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Host: fmt.Sprintf("%s:%s", p.Ip, p.Port)})},
		Timeout:   timeout}

	var testUrl string
	if https {
		testUrl = "https://ip.cip.cc"
	} else {
		testUrl = "http://ip.cip.cc"
	}

	resp, err := client.Get(testUrl)
	if err != nil {
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	if 200 != resp.StatusCode {
		return errors.New(fmt.Sprintf("http code %d", resp.StatusCode))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html := strings.TrimSpace(string(b))
	if html != p.GetIp() {
		return errors.New(fmt.Sprintf("error resp"))
	}
	latency := time.Now().UnixNano() - startAt.UnixNano()
	p.Latency = int(latency / 1000 / 1000)
	if https {
		p.Schema = "https"
	} else {
		p.Schema = "http"
	}
	return
}
