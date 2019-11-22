package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"time"
	//"strings"
)

func ServeReverse() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.ProxyBind, config.ProxyPort))
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
	for i := 1; i < config.Retry+1; i++ {

		p := proxies[rand.Intn(l)]
		serverAddr := p.GetProxyUrl()
		logger.WithField("time", i).WithField("serverAddr", serverAddr).Debug("try connect server")
		server, err = tryExchange(serverAddr, schema, connString)
		if err != nil {
			if !os.IsTimeout(err) {
				logger.WithField("serverAddr", serverAddr).WithError(err).Error("fatal error establish connect")
			} else {
				logger.WithField("serverAddr", serverAddr).WithError(err).Error("timeout")
			}
			continue
		}

		break

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
		//_, _ = fmt.Fprint(server, connString+"\r\n")
		_, err = server.Write(b[:n])
		if err != nil {
			logger.WithError(err).Error("Write")
		}
	}

	go io.Copy(server, client)
	_, err = io.Copy(client, server)
	if err != nil {
		logger.WithError(err).Error("io.Copy")
	}
	logger.Debug("closed")
	return
}

func tryExchange(serverAddr, schema, connString string) (server net.Conn, err error) {
	server, err = net.DialTimeout("tcp", serverAddr, time.Duration(config.TcpTimeout)*time.Second)
	if err != nil {
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
