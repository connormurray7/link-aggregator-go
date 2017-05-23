package linkagg

import (
	"fmt"
	"net/http"
	"time"

	"io/ioutil"

	"log"

	"github.com/spf13/viper"
)

//Handler makes a request to the outside APIs
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

//Server implements the requester interface and calls out to external APIs.
type Server struct {
	cache        *Cache
	config       *viper.Viper
	client       *http.Client
	reqTimes     chan int64
	maxReqPerSec int
}

//NewServer constructs a new link agg server given a config file.
func NewServer(config *viper.Viper) Server {
	var server Server

	server.config = config
	server.cache = NewLinkAggCache(config)
	server.client = &http.Client{
		Timeout: time.Second * 10,
	}
	server.maxReqPerSec = config.GetInt("ratelimit")
	server.reqTimes = make(chan int64, server.maxReqPerSec)
	return server
}

//Handle fetches all of the information from external APIs if not cached.
func (server *Server) Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request")
	defer r.Body.Close()
	arr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse request.", r)
	}
	req := string(arr)
	result := server.cache.Get(req)
	if result == "" && !server.needRateLimit() {
		resp := FetchExternalRequest(req, server.config, server.client)
		server.cache.Set(req, resp)
		w.Write([]byte(resp))
	}
}

func (server *Server) needRateLimit() bool {
	if len(server.reqTimes) < server.maxReqPerSec {
		return false
	}
	curTime := time.Now().Unix()
	server.reqTimes <- curTime
	time := <-server.reqTimes
	return int(curTime-time) < server.maxReqPerSec
}
