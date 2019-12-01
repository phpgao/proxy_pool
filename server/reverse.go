package server

import (
	"crypto/tls"
	"fmt"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

var Server *http.Server

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	proxies, err := storeEngine.Get(map[string]string{
		"schema": "https",
	})
	if err != nil {
		return
	}
	l := len(proxies)
	if l == 0 {
		http.Error(w, "no valid proxy", http.StatusServiceUnavailable)
		return
	}
	proxy := proxies[rand.Intn(l)]

	msg := fmt.Sprintf(model.ConnectCommand, http.MethodConnect, r.Host, "HTTP/1.1", r.Host)

	destConn, err := net.DialTimeout("tcp", proxy.GetProxyUrl(), time.Duration(util.ServerConf.HttpsConnectTimeOut)*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	_, err = destConn.Write([]byte(msg))

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		destConn.Close()
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
func handleHTTP(w http.ResponseWriter, req *http.Request) {
	proxy, err := storeEngine.Random()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	Transport := http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Host: proxy.GetProxyUrl()},
		),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		// skip cert check
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := Transport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func ServeReverse() {
	addr := fmt.Sprintf("%s:%d", util.ServerConf.ProxyBind, util.ServerConf.ProxyPort)
	Server = &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	logger.WithField("addr", addr).Info("dynamic proxy listen and serve")

	Server.ListenAndServe()
}
