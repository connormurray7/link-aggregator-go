package linkagg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//Requester accepts a string query and returns an answer.
type Requester interface {
	Request(query string) string
}

//RequestService makes http requests to list of external APIs.
type RequestService struct {
	apis   []*ExternalAPI
	client *http.Client
}

//Message contains the information for every row in a response.
type Message struct {
	Title string
	Link  string
}

//NewRequestService requires Viper config for API information.
func NewRequestService(config *viper.Viper) *RequestService {
	var r RequestService

	r.apis = []*ExternalAPI{
		NewGithubAPI(config),
		NewHackerNewsAPI(config),
		NewStackOverflowAPI(config),
	}
	r.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return &r
}

//Request calls all external apis and returns a json string of parsed responses.
func (r *RequestService) Request(query string) string {
	m := make(map[string]*[]Message)
	for _, api := range r.apis {
		m[api.name] = api.makeExternalRequest(query, r.client)
	}
	result, err := json.Marshal(m)
	if err != nil {
		log.Print("Unable to encode response", err)
	}
	return string(result)
}

func (e *ExternalAPI) makeExternalRequest(query string, client *http.Client) *[]Message {
	req, err := http.NewRequest("GET", e.url, nil)
	if err != nil {
		log.Print("Error creating new request", err)
		return nil
	}
	q := req.URL.Query()
	for _, param := range e.params {
		q.Add(param.Key, param.Value)

	}
	q.Add(e.queryKey, query)
	req.URL.RawQuery = q.Encode()
	json := executeRequest(req, client)
	return parseJSONResponse(json, e.parsing)
}

func executeRequest(req *http.Request, client *http.Client) string {
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

func parseJSONResponse(json string, params ParsingParams) *[]Message {
	result := gjson.Get(json, params.items)
	parsed := make([]Message, 20)
	num := 0

	for _, hit := range result.Array() {
		record := hit.Map()
		parsed[num] = Message{record[params.title].Str, record[params.url].Str}
		num++
	}
	return &parsed
}
