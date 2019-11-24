package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/source"
	"github.com/phpgao/proxy_pool/util"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

var (
	config         = util.GetConfig()
	logger         = util.GetLogger()
	storeEngine    = db.GetDb()
	IdleConnClosed = make(chan struct{})
	Srv            *http.Server
)

var home = "https://github.com/phpgao/proxy_pool"

type Resp struct {
	Code   int         `json:"code"`
	Error  string      `json:"error"`
	Total  int         `json:"total"`
	Schema interface{} `json:",omitempty"`
	Score  interface{} `json:",omitempty"`
	Cn     int         `json:",omitempty"`
	Data   interface{} `json:"data"`
	Get    string      `json:",omitempty"`
	Random string      `json:",omitempty"`
	Home   string      `json:",omitempty"`
}

func Serve() {
	addr := fmt.Sprintf("%s:%d", config.ApiBind, config.ApiPort)

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
	mux.HandleFunc("/random", handlerRandom)
	mux.HandleFunc("/get", handlerQuery)
	mux.HandleFunc("/", handlerStatus)
	return mux
}

func handlerQuery(w http.ResponseWriter, r *http.Request) {
	resp := Resp{
		Code: http.StatusOK,
	}
	proxies, err := Filter(r, resp)
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Data = proxies
		resp.Total = len(proxies)
	}

	respText, err := json.Marshal(resp)
	if err != nil {
		resp.Error = err.Error()
	}
	_, _ = w.Write(respText)
}

func handlerStatus(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	resp := Resp{
		Code: http.StatusOK,
	}

	allSpider := source.ListOfSpider

	status := make(map[string]int)
	for _, s := range allSpider {
		status[s.Name()] = 0
	}
	scores := make(map[string]int)
	schema := make(map[string]int)
	country := 0

	proxies := storeEngine.GetAll()
	l := len(proxies)
	if l > 0 {
		resp.Total = len(proxies)
		for _, p := range proxies {
			if p.Schema == "http" {
				setDefault(schema, "http", 0, 1)
			} else {
				setDefault(schema, "https", 0, 1)
			}
			status[p.From] += 1
			setDefault(scores, strconv.Itoa(p.Score), 0, 1)
			if p.Country == "cn" {
				country += 1
			}
		}
	}
	resp.Schema = schema
	resp.Data = status
	resp.Score = scores
	resp.Cn = country
	resp.Home = home
	resp.Get = "/get?schema=&score="
	resp.Random = "/random?schema=&score="
	respText, err := json.Marshal(resp)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(respText)
}

func handlerRandom(w http.ResponseWriter, r *http.Request) {
	resp := Resp{
		Code: http.StatusOK,
	}
	proxies, err := Filter(r, resp)
	if err != nil {
		resp.Error = err.Error()
	} else {
		if len(proxies) > 0 {
			resp.Data = proxies[rand.Intn(len(proxies))]
			resp.Total = len(proxies)
		} else {
			resp.Data = nil
			resp.Total = 0
		}
	}

	respText, err := json.Marshal(resp)
	if err != nil {
		resp.Error = err.Error()
	}
	_, _ = w.Write(respText)
}

func Filter(r *http.Request, resp Resp) (proxies []model.HttpProxy, err error) {
	err = r.ParseForm()
	// http or https,default all
	schema := r.FormValue("schema")
	// ip in China or not
	// "1" -> China only
	// else
	cn := r.FormValue("cn")
	// score above given number
	score := r.FormValue("score")
	_source := r.FormValue("source")
	country := r.FormValue("country")
	if err != nil {
		resp.Error = err.Error()
	}
	proxies, err = storeEngine.Get(map[string]string{
		"schema":  schema,
		"cn":      cn,
		"score":   score,
		"source":  _source,
		"country": country,
	})
	return
}

func setDefault(h map[string]int, k string, v, inc int) (set bool, r int) {
	if _, set = h[k]; !set {
		h[k] = v
		set = true
	}
	h[k] += inc
	return
}
