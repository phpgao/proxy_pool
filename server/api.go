package server

import (
	"github.com/gin-gonic/gin"
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/job"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"math/rand"
	"net/http"
	"strconv"
)

var (
	logger      = util.GetLogger("reverse")
	storeEngine = db.GetDb()
	//IdleConnClosed = make(chan struct{})
	//Srv            *http.Server
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

//
//func Serve() {
//
//	var err error
//	addr := fmt.Sprintf("%s:%d", util.ServerConf.ApiBind, util.ServerConf.ApiPort)
//
//	logger.WithField("addr", addr).Info("api listen and serve")
//	Srv = &http.Server{Addr: addr, Handler: GetMux()}
//
//	go func() {
//		sigint := make(chan os.Signal, 1)
//		signal.Notify(sigint, os.Interrupt)
//		<-sigint
//		logger.Info("shutting down server")
//		if err = Srv.Shutdown(context.Background()); err != nil {
//			logger.WithError(err).Error("api server Shutdown")
//		}
//		if err = Server.Shutdown(context.Background()); err != nil {
//			logger.WithError(err).Error("proxy server Shutdown")
//		}
//		close(IdleConnClosed)
//	}()
//
//	if err = Srv.ListenAndServe(); err != http.ErrServerClosed {
//		logger.Fatalf("HTTP server ListenAndServe: %v", err)
//	}
//
//	<-IdleConnClosed
//}

//func GetMux() *http.ServeMux {
//	mux := http.NewServeMux()
//	mux.HandleFunc("/random", handlerRandom)
//	mux.HandleFunc("/get", handlerQuery)
//	mux.HandleFunc("/", handlerStatus)
//	return mux
//}

func routerApi() http.Handler {
	if !util.ServerConf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.Default()
	e.GET("/", handlerStatus)
	e.GET("/get", handlerQuery)
	e.GET("/random", handlerRandom)

	return e
}

func handlerQuery(c *gin.Context) {
	var err error
	resp := Resp{
		Code: http.StatusOK,
	}

	proxies, err := Filter(c)

	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Data = proxies
		resp.Total = len(proxies)
	}

	c.JSON(http.StatusOK, resp)
}

func handlerStatus(c *gin.Context) {
	resp := Resp{
		Code: http.StatusOK,
	}
	allSpider := job.ListOfSpider
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
	c.JSON(http.StatusOK, resp)
}

func handlerRandom(c *gin.Context) {
	resp := Resp{
		Code: http.StatusOK,
	}
	proxies, err := Filter(c)
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

	c.JSON(http.StatusOK, resp)
}

func Filter(c *gin.Context) (proxies []model.HttpProxy, err error) {

	// http or https,default all
	schema := c.Query("schema")
	// ip in China or not
	// "1" -> China only
	// else
	cn := c.Query("cn")
	// score above given number
	score := c.Query("score")
	_source := c.Query("source")
	country := c.Query("country")

	return storeEngine.Get(map[string]string{
		"schema":  schema,
		"cn":      cn,
		"score":   score,
		"source":  _source,
		"country": country,
	})
}

func setDefault(h map[string]int, k string, v, inc int) (set bool, r int) {
	if _, set = h[k]; !set {
		h[k] = v
		set = true
	}
	h[k] += inc
	return
}
