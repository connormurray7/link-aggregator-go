package linkagg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//Message contains the information for every row in a response.
type Message struct {
	Title string
	Link  string
}

// type Requester interface {
// 	MakeExternalRequest(url string, params []EncodingPair, client *http.Client) string
// }

type ExternalApi struct {
	params   []EncodingPair
	queryKey string
	url      string
	client   *http.Client
}

type EncodingPair struct {
	Key   string
	Value string
}

//FetchExternalRequest calls all external apis and returns a json string of parsed responses.
func FetchExternalRequest(query string, config *viper.Viper, client *http.Client) string {
	m := make(map[string]*[]Message)
	m["Github"] = makeGithubRequest(query, config, client)
	m["Hacker News"] = makeHackerNewsRequest(query, config, client)
	m["Stack Overflow"] = makeStackOverflowRequest(query, config, client)
	result, err := json.Marshal(m)
	if err != nil {
		log.Print("Unable to encode response", err)
	}
	return string(result)
}

func (e *ExternalApi) MakeExternalRequest(query string) string {
	req, err := http.NewRequest("GET", e.url, nil)
	if err != nil {
		log.Print("Error creating new Github request", err)
		return ""
	}
	q := req.URL.Query()
	for _, param := range e.params {
		q.Add(param.Key, param.Value)

	}
	q.Add(e.queryKey, query)
	req.URL.RawQuery = q.Encode()
	json := makeRequest(req, e.client)
	return json
}

func makeRequest(req *http.Request, client *http.Client) string {
	log.Println("Making request", req.URL.String())
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

func parseJSONResponse(json string, items string, title string, url string) *[]Message {
	result := gjson.Get(json, items)
	parsed := make([]Message, 20)
	num := 0

	for _, hit := range result.Array() {
		record := hit.Map()
		parsed[num] = Message{record[title].Str, record[url].Str}
		num++
	}
	return &parsed
}
