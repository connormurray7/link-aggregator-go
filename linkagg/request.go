package linkagg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"net/http/httputil"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//Message contains the information for every row in a response.
type Message struct {
	title string
	link  string
}

//FetchExternalRequest calls all external apis and returns a json string of parsed responses.
func FetchExternalRequest(query string, config *viper.Viper, client *http.Client) string {
	m := make(map[string][]Message)
	m["Github"] = makeGithubRequest(query, config, client)
	m["Hacker News"] = makeHackerNewsRequest(query, config, client)
	m["Stack Overflow"] = makeStackOverflowRequest(query, config, client)
	result, err := json.Marshal(m)
	if err != nil {
		log.Print("Unable to encode response", err)
	}
	return string(result)
}

func makeHackerNewsRequest(query string, config *viper.Viper, client *http.Client) []Message {
	req, err := http.NewRequest("GET", config.GetString("HackerNews.url"), nil)
	if err != nil {
		log.Print("Error creating new Hacker News request", err)
		return nil
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("tags", "story")
	q.Set("hitsPerPage", "15")

	json := makeRequest(req, client)
	log.Println("Received request from Hacker News", json)
	return parseJSONResponse(json, "hits", "title", "url")
}

func makeStackOverflowRequest(query string, config *viper.Viper, client *http.Client) []Message {
	req, err := http.NewRequest("GET", config.GetString("StackOverflow.url"), nil)
	if err != nil {
		log.Print("Error creating new Stack Overflow request", err)
		return nil
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("order", "desc")
	q.Set("sort", "relevance")
	q.Set("accepted", "True")
	q.Set("site", "stackoverflow")
	q.Set("pagesize", "15")

	json := makeRequest(req, client)
	log.Println("Received request from Stack Overflow", json)
	return parseJSONResponse(json, "items", "title", "link")
}

func makeGithubRequest(query string, config *viper.Viper, client *http.Client) []Message {
	req, err := http.NewRequest("GET", config.GetString("Github.url"), nil)
	if err != nil {
		log.Print("Error creating new Github request", err)
		return nil
	}
	q := req.URL.Query()
	q.Set("q", query)
	q.Set("sort", "stars")
	q.Set("per_page", "15")

	json := makeRequest(req, client)
	log.Println("Received request from Github", json)
	return parseJSONResponse(json, "items", "name", "html_url")
}

func makeRequest(req *http.Request, client *http.Client) string {
	reqStr, err := httputil.DumpRequest(req, false)
	if err != nil {
		log.Println("Malformed request", req)
	}
	log.Println("Making request", string(reqStr))
	resp, err := client.Do(req)
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
