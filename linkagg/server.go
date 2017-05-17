package linkagg

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

//Handler makes a request to the outside APIs
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

//Server implements the requester interface and calls out to external APIs.
type Server struct {
	cache  *Cache
	config *viper.Viper
	client *http.Client
}

//NewServer constructs a new link agg server given a config file.
func NewServer(config *viper.Viper) Server {
	var server Server

	server.config = config
	server.cache = NewLinkAggCache(config)
	server.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return server
}

//Handle fetches all of the information from external APIs if not cached.
func (server *Server) Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request")
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	req := buf.String()
	result := server.cache.Get(req)
	if result == "" {
		resp := FetchExternalRequest(req, server.config, server.client)
		server.cache.Set(req, resp)
		w.Write([]byte(resp))
	}
}
