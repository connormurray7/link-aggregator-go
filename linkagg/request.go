package linkagg

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//Requester makes a request to the outside APIs
type Requester interface {
	Request(req string) string
}

//LinkAgg implements the requester interface and calls out to external APIs.
type LinkAgg struct {
	cache  *Cache
	config *viper.Viper
	client *http.Client
}

type LinkAggMessage struct {
	title string
	link  string
}

//NewLinkAgg constructs a link agg given a config file.
func NewLinkAgg(config *viper.Viper) LinkAgg {
	var linkAgg LinkAgg

	linkAgg.config = config
	linkAgg.cache = NewLinkAggCache(config)
	linkAgg.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return linkAgg
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
	hnRequest := linkAgg.makeHackerNewsRequest(query)
	soRequest := linkAgg.makeStackOverflowRequest(query)
	ghRequest := linkAgg.makeGithubRequest(query)
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
	q.Set("hitsPerPage", "15")

	resp := linkAgg.makeRequest(req)

	return ""
}

func (linkAgg *LinkAgg) makeStackOverflowRequest(query string) string {
	req, err := http.NewRequest("GET", linkAgg.config.GetString("StackOverflow.url"), nil)
	if err != nil {
		log.Print(err)
		return ""
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("order", "desc")
	q.Set("accepted", "True")
	q.Set("site", "stackoverflow")
	q.Set("pagesize", "15")

	resp := linkAgg.makeRequest(req)
	parsed := []LinkAggMessage{}

	for hit := range resp["items"] {
		log.Print(hit["title"], hit["link"])
	}

	return ""
}

func (linkAgg *LinkAgg) makeGithubRequest(query string) []LinkAggMessage {
	req, err := http.NewRequest("GET", linkAgg.config.GetString("Github.url"), nil)
	if err != nil {
		log.Print(err)
		return ""
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("sort", "stars")
	q.Set("per_page", "15")

	json := linkAgg.makeRequest(req)
	result := gjson.Get(json, "items")
	parsed := make([]LinkAggMessage, 20)
	num := 0

	for _, hit := range result.Array() {
		record := hit.Map()
		log.Print(record["title"], record["link"])
		parsed[num] = LinkAggMessage{record["title"].String(), record["link"].String()}
		num++
	}
	return parsed[:num]
}

func (linkAgg *LinkAgg) makeRequest(req *http.Request) string {
	resp, err := linkAgg.client.Do(req)
	if err != nil {
		log.Print("Unable to complete request", err)
		return ""
	}
	byteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Unable to read all bytes in response", resp, err)
		return ""
	}

	return string(byteArr[:])
}
