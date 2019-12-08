package server

import (
	"crypto/tls"
	"fmt"
	"github.com/phpgao/proxy_pool/cache"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	Server  *http.Server
	timeOut = time.Duration(util.ServerConf.HttpsConnectTimeOut) * time.Second
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	var err error
	p := *cache.Cache.Get()
	proxies := p["https"]
	var destConn net.Conn
	if proxies == nil {
		logger.Debug("serve as a https proxy")
		destConn, err = net.DialTimeout("tcp", r.Host, timeOut)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)

	} else {
		l := len(proxies)
		proxy := proxies[rand.Intn(l)]
		logger.WithField("proxy", proxy.GetProxyWithSchema()).Debug("dynamic https")
		msg := fmt.Sprintf(model.ConnectCommand, http.MethodConnect, r.Host, "HTTP/1.1", r.Host)

		destConn, err = net.DialTimeout("tcp", proxy.GetProxyUrl(), timeOut)
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

		err = destConn.SetReadDeadline(time.Now().Add(timeOut))

		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			destConn.Close()
			return
		}

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

	defer func() {
		err := destination.Close()
		if err != nil {
			logger.WithError(err).Warn("error close remote conn")
		}
		source.Close()
	}()
	io.Copy(destination, source)
}
func handleHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	var Transport http.RoundTripper

	p := *cache.Cache.Get()
	proxies := p["http"]
	if proxies == nil {
		logger.Debug("serve as a http proxy")
		Transport = http.DefaultTransport
	} else {
		proxy := proxies[rand.Intn(len(proxies))]
		logger.WithField("proxy", proxy.GetProxyWithSchema()).Debug("dynamic http")
		Transport = &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Host: proxy.GetProxyUrl()},
			),
			DialContext: (&net.Dialer{
				Timeout:   timeOut,
				KeepAlive: timeOut,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       timeOut,
			TLSHandshakeTimeout:   timeOut,
			ExpectContinueTimeout: 1 * time.Second,
			// skip cert check
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
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

	if err := Server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
