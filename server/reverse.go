package server

import (
	"crypto/tls"
	"fmt"
	"github.com/phpgao/proxy_pool/model"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

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

	destConn, err := net.DialTimeout("tcp", proxy.GetProxyUrl(), 10*time.Second)
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
	//proxies, err := storeEngine.Get(map[string]string{
	//	"schema": "http",
	//})
	//if err != nil {
	//	return
	//}
	//l := len(proxies)
	//if l == 0 {
	//	http.Error(w, "no valid proxy", http.StatusServiceUnavailable)
	//	return
	//}
	//proxy := proxies[rand.Intn(l)]

	proxy, err := storeEngine.Random()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	Transport := http.Transport{
		Proxy: func(request *http.Request) (url *url.URL, err error) {
			proxyURL, err := url.Parse("http://" + proxy.GetProxyUrl())
			if err != nil {
				return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
			}
			return proxyURL, nil
		},
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
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.ProxyBind, config.ProxyPort),
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
	server.ListenAndServe()
}
