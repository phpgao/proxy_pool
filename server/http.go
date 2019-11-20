package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"net/http"
	"os"
	"os/signal"
)

var (
	config         = util.GetConfig()
	logger         = util.GetLogger()
	storeEngine    = db.GetDb()
	IdleConnClosed = make(chan struct{})
	Srv            *http.Server
)

type JsonResp struct {
	Code  int               `json:"code"`
	Error string            `json:"error"`
	Total int               `json:"total"`
	Data  []model.HttpProxy `json:"data"`
}

func Serve() {
	addr := fmt.Sprintf(":%d", config.Port)
	logger.WithField("addr", addr).Info("listen and serve")
	Srv = &http.Server{Addr: addr, Handler: GetMux()}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		logger.Info("shutting down server")
		if err := Srv.Shutdown(context.Background()); err != nil {
			logger.WithError(err).Error("HTTP server Shutdown")
		}
		close(IdleConnClosed)
	}()

	if err := Srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-IdleConnClosed
}

func GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/all", handlerAll)
	mux.HandleFunc("/random", handlerRandom)
	return mux
}

func handlerAll(w http.ResponseWriter, _ *http.Request) {
	proxies := storeEngine.GetAll()
	resp := JsonResp{
		Code:  200,
		Error: "",
		Total: len(proxies),
		Data:  proxies,
	}
	respText, err := json.Marshal(resp)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(respText)
}

func handlerRandom(w http.ResponseWriter, _ *http.Request) {
	randomProxy, err := storeEngine.Random()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(randomProxy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}
