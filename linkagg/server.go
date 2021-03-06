package linkagg

import (
	"net/http"
	"time"

	"io/ioutil"

	"log"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//Server implements the requester interface and calls out to external APIs.
type Server struct {
	cache        *Cache
	config       *viper.Viper
	reqTimes     chan int64
	maxReqPerSec int
	service      *RequestService
}

//NewServer constructs a new link agg server given a config file.
func NewServer(config *viper.Viper) *Server {
	var server Server

	server.config = config
	server.cache = NewLinkAggCache(config)
	server.maxReqPerSec = config.GetInt("ratelimit")
	server.reqTimes = make(chan int64, server.maxReqPerSec)
	server.service = NewRequestService(config)

	return &server
}

//Handle fetches all of the information from external APIs if not cached.
func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	arr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to parse request.", r)
	}
	jsonReq := string(arr)
	req := gjson.Get(jsonReq, "term").String()
	log.Println("Inbound request", req)
	result := server.cache.Get(req)
	if result == "" && !server.needRateLimit() {
		log.Println("Request not cached, fetching from external APIs")
		result = server.service.Request(req)
		server.cache.Set(req, result)
	}
	log.Println("Sending back: ", result)
	w.Write([]byte(result))
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
