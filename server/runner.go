package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/phpgao/proxy_pool/util"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	g              errgroup.Group
	IdleConnClosed = make(chan struct{})
)

func RunService() {
	var ApiService, ProxyService *http.Server

	if util.ServerConf.EnableApi {
		addr := fmt.Sprintf("%s:%d", util.ServerConf.ApiBind, util.ServerConf.ApiPort)
		ApiService = &http.Server{
			Addr:         addr,
			Handler:      routerApi(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		ApiService.SetKeepAlivesEnabled(false)

		g.Go(func() error {
			return ApiService.ListenAndServe()
		})

	}

	if util.ServerConf.EnableProxy {
		addr := fmt.Sprintf("%s:%d", util.ServerConf.ProxyBind, util.ServerConf.ProxyPort)
		ProxyService = &http.Server{
			Addr:         addr,
			Handler:      routerProxy(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}
		ProxyService.SetKeepAlivesEnabled(false)

		g.Go(func() error {
			logger.WithField("addr", addr).Info("ProxyService listen and serve")
			return ProxyService.ListenAndServe()
		})
	}

	go func() {
		var err error
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		logger.Info("shutting down server")
		if ApiService != nil {
			logger.Info("shutting down api server")
			if err = ApiService.Shutdown(context.Background()); err != nil {
				logger.WithError(err).Error("Could not gracefully shutdown api server")
			}
		}
		if ProxyService != nil {
			logger.Info("shutting down proxy server")
			if err = ProxyService.Shutdown(context.Background()); err != nil {
				logger.WithError(err).Error("Could not gracefully shutdown proxy server")
			}
		}
		close(IdleConnClosed)
	}()

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	<-IdleConnClosed
	os.Exit(0)
}
