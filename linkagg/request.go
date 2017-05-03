package linkagg

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

//Requester makes a request to the outside APIs
type Requester interface {
	Request(req string) string
}

//LinkAgg implements the requester interface and calls out to external APIs.
type LinkAgg struct {
	cache  *Cache
	config *viper.Viper
}

//Request fetches all of the information from external APIs if not cached.
func (linkAgg *LinkAgg) Request(req string) string {
	result := linkAgg.cache.Get(req)
	if result != "" {
		return result
	}
	resp := linkAgg.fetchExternalRequest(req)
	linkAgg.cache.Set(req, resp)
	return resp
}

func (linkAgg *LinkAgg) fetchExternalRequest(query string) string {
	hnRequest := linkAgg.makeHackerNewsRequest()
	soRequest := linkAgg.createStackOverflowRequest()
	ghRequest := linkAgg.createGithubRequest()
}

func (linkAgg *LinkAgg) makeHackerNewsRequest(query string) string {
	req, err := http.NewRequest("GET", linkAgg.config.GetString("HackerNews.url"), nil)
	if err != nil {
		log.Print(err)
		return ""
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("tags", "story")
	q.Set("hitsPerPage", 15)

}

func (linkAgg *LinkAgg) createStackOverflowRequest() string {

}

func (linkAgg *LinkAgg) createGithubRequest() string {

}
