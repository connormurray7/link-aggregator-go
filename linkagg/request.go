package linkagg

import (
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
	hnRequest := createHackerNewsRequest()
	soRequest := createStackOverflowRequest()
	ghRequest := createGithubRequest()
}

func createHackerNewsRequest() string {

}

func createStackOverflowRequest() string {

}

func createGithubRequest() string {

}
