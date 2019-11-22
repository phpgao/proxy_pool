package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/pool"
	"github.com/phpgao/proxy_pool/ppool"
	"io"
	"log"
	//"math/rand"
	"net"
	"net/url"
	"os"
	"time"
)

func ServeReverse() {
	addr := fmt.Sprintf("%s:%d", config.ProxyBind, config.ProxyPort)
	logger.WithField("addr", addr).Info("listen and serve")

	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.WithError(err).Fatal("ServeReverse")
	}
	for {
		client, err := l.Accept()
		if err != nil {
			logger.WithError(err).Error("ServeReverse Accept")
			continue
		}

		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var connString = string(b[:bytes.IndexByte(b[:], '\r')])
	var method, host, version string
	_, _ = fmt.Sscanf(connString, "%s %s %s", &method, &host, &version)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	var schema = "http"
	if hostPortURL.Opaque == "443" {
		schema = "https"
		connString = fmt.Sprintf("%s %s %s\r\n\r\n", method, host, version)
	}

	proxies, err := storeEngine.Get(map[string]string{
		"schema": schema,
	})
	if err != nil {
		logger.WithError(err).Error("err get proxies")
		return
	}
	l := len(proxies)
	if l == 0 {
		logger.Error("no available proxy")
		return
	}
	var server net.Conn
	server, err = tryExchange(schema, connString)
	if err != nil {
		if !os.IsTimeout(err) {
			logger.WithError(err).Error("fatal error establish connect")
		} else {
			logger.WithError(err).Error("timeout")
		}
		return
	}
	if server == nil {
		return
	}
	defer server.Close()
	if method == "CONNECT" {
		_, err = fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
		if err != nil {
			logger.WithError(err).Error("CONNECT")
		}
	} else {
		_, err = server.Write(b[:n])
		if err != nil {
			logger.WithError(err).Error("Write")
			if pc, ok := server.(*pool.PoolConn); ok {
				pc.MarkUnusable()
				pc.Close()
			}
			return
		}
	}

	go func() {
		defer server.Close()
		defer client.Close()
		buf := make([]byte, 2048)
		_, err := io.CopyBuffer(server, client, buf)
		if err != nil {

			logger.WithError(err).Error("error io.Copy goroutine")
		}
		if pc, ok := server.(*pool.PoolConn); ok {
			pc.MarkUnusable()
			logger.Warn("MarkUnusable goroutine")

			pc.Close()
		}
	}()
	buf := make([]byte, 2048)
	_, err = io.CopyBuffer(client, server, buf)
	if err != nil {
		logger.WithError(err).Error("error io.Copy outside")
	}
	if pc, ok := server.(*pool.PoolConn); ok {
		pc.MarkUnusable()
		logger.Warn("MarkUnusable")

		pc.Close()
	}
	logger.WithField("chan", ppool.Http.Len()).Error("len")
	return
}

func tryExchange(schema, connString string) (server net.Conn, err error) {
	server, err = ppool.Http.Get()
	if err != nil {
		return
	}

	if server == nil {
		logger.WithField("server.RemoteAddr()", server).WithField("err", err).Debug("server")
		err = errors.New("nil server")
		return
	}

	if schema == "https" {
		// need to send connect again
		_, err = fmt.Fprint(server, connString)
		if err != nil {
			return
		}
		// read 200 code
		var mb [1024]byte
		if err = server.SetReadDeadline(time.Now().Add(time.Duration(config.ProxyTimeout) * time.Second)); err != nil {
			return
		}
		_, err = server.Read(mb[:])
		if err != nil {

			return
		}

		var stringBack = string(mb[:bytes.IndexByte(mb[:], '\r')])
		var code, version string

		_, err = fmt.Sscanf(stringBack, "%s %s", &version, &code)

		if err != nil {
			return
		}

		if version != "HTTP/1.1" || code != "200" {
			return nil, errors.New(stringBack)
		}
	}
	return
}
