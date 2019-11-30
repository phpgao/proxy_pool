package model

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/phpgao/proxy_pool/util"
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

var (
	config = util.ServerConf
	logger = util.GetLogger()
)

const ConnectCommand = "%s %s %s\r\nHost: %s\r\nProxy-Connection: Keep-Alive\r\n\r\n"

type HttpProxy struct {
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Schema    string `json:"schema"`
	Score     int    `json:"score"`
	Latency   int    `json:"latency"`
	From      string `json:"-"`
	Anonymous int    `json:"anonymous"`
	Country   string `json:"country"`
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

func (p *HttpProxy) GetProxyWithSchema() string {
	return fmt.Sprintf("%s://%s:%s", p.Schema, p.Ip, p.Port)
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
func (p *HttpProxy) IsHttps() bool {
	return p.Schema == "https"
}

func (p *HttpProxy) SimpleTcpTest(timeOut time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", p.Ip, p.Port), timeOut)
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
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Host: fmt.Sprintf("%s:%s", p.Ip, p.Port)},
			),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout}

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

// test tcp
func (p *HttpProxy) TestTcp() (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%s", p.Ip, p.Port), config.GetTcpTestTimeOut())
	if err != nil {
		if conn != nil {
			_ = conn.Close()
		}
		return
	}

	return
}

// test http connect method
func (p *HttpProxy) TestConnectMethod(conn net.Conn) (err error) {
	defer conn.Close()
	testHost := "cip.cc:443"
	Connect := fmt.Sprintf(ConnectCommand, http.MethodConnect, testHost, "HTTP/1.1", testHost)
	_, err = conn.Write([]byte(Connect))
	if err != nil {
		return
	}
	// read 200 code
	var mb [1024]byte
	if err = conn.SetReadDeadline(time.Now().Add(time.Duration(config.ProxyTimeout) * time.Second)); err != nil {
		return
	}
	_, err = conn.Read(mb[:])
	if err != nil {
		return
	}
	firstLineIndex := bytes.IndexByte(mb[:], '\n')
	if firstLineIndex == -1 {
		return errors.New("error response format")
	}
	var stringBack = string(mb[:firstLineIndex])
	var code, version string

	_, err = fmt.Sscanf(stringBack, "%s %s", &version, &code)

	if err != nil {
		return
	}

	if version != "HTTP/1.1" || code != "200" {
		return errors.New("bad response" + stringBack)
	}

	return
}
