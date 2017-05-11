package linkagg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//Handler makes a request to the outside APIs
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

//LinkAgg implements the requester interface and calls out to external APIs.
type LinkAgg struct {
	cache  *Cache
	config *viper.Viper
	client *http.Client
}

//Message contains the information for every row in a response.
type Message struct {
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

//Handle fetches all of the information from external APIs if not cached.
func (linkAgg *LinkAgg) Handle(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	req := buf.String()
	result := linkAgg.cache.Get(req)
	if result == "" {
		resp := linkAgg.fetchExternalRequest(req)
		linkAgg.cache.Set(req, resp)
		w.Write([]byte(resp))
	}
}

func (linkAgg *LinkAgg) fetchExternalRequest(query string) string {
	m := make(map[string][]Message)
	m["Github"] = linkAgg.makeGithubRequest(query)
	m["Hacker News"] = linkAgg.makeHackerNewsRequest(query)
	m["Stack Overflow"] = linkAgg.makeStackOverflowRequest(query)
	result, err := json.Marshal(m)
	if err != nil {
		log.Print("Unable to encode response", err)
	}
	return string(result)
}

func (linkAgg *LinkAgg) makeHackerNewsRequest(query string) []Message {
	req, err := http.NewRequest("GET", linkAgg.config.GetString("HackerNews.url"), nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("tags", "story")
	q.Set("hitsPerPage", "15")

	json := linkAgg.makeRequest(req)
	return parseJSONResponse(json, "hits", "title", "url")
}

func (linkAgg *LinkAgg) makeStackOverflowRequest(query string) []Message {
	req, err := http.NewRequest("GET", linkAgg.config.GetString("StackOverflow.url"), nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("order", "desc")
	q.Set("accepted", "True")
	q.Set("site", "stackoverflow")
	q.Set("pagesize", "15")

	json := linkAgg.makeRequest(req)
	return parseJSONResponse(json, "items", "title", "link")
}

func (linkAgg *LinkAgg) makeGithubRequest(query string) []Message {
	req, err := http.NewRequest("GET", linkAgg.config.GetString("Github.url"), nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("sort", "stars")
	q.Set("per_page", "15")

	json := linkAgg.makeRequest(req)
	return parseJSONResponse(json, "items", "name", "html_url")
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

func parseJSONResponse(json string, items string, title string, url string) []Message {
	result := gjson.Get(json, items)
	parsed := make([]Message, 20)
	num := 0

	for _, hit := range result.Array() {
		record := hit.Map()
		parsed[num] = Message{record[title].String(), record[url].String()}
		num++
	}
	return parsed[:num]
}
