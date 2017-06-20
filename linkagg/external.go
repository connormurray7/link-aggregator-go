package linkagg

import "github.com/spf13/viper"

//ExternalAPI contains all specific info to 3rd party APIs.
type ExternalAPI struct {
	name     string
	queryKey string
	url      string
	params   []URLParam
	parsing  ParsingParams
}

//URLParam is a key value pair.
type URLParam struct {
	Key   string
	Value string
}

//ParsingParams contains the names of keys in JSON response.
type ParsingParams struct {
	items string
	title string
	url   string
}

//NewGithubApi requires a config to for the url.
func NewGithubAPI(config *viper.Viper) *ExternalAPI {
	var e ExternalAPI

	e.name = "Github"
	e.url = config.GetString("Github.url")
	e.queryKey = "q"
	e.params = []URLParam{
		URLParam{"sort", "stars"},
		URLParam{"per_page", "15"},
	}
	e.parsing = ParsingParams{
		items: "items",
		title: "name",
		url:   "html_url",
	}
	return &e
}

//NewHackerNewsAPI requires a config to for the url.
func NewHackerNewsAPI(config *viper.Viper) *ExternalAPI {
	var e ExternalAPI

	e.name = "Hacker News"
	e.url = config.GetString("HackerNews.url")
	e.queryKey = "query"
	e.params = []URLParam{
		URLParam{"tags", "story"},
		URLParam{"hitsPerPage", "15"},
	}
	e.parsing = ParsingParams{
		items: "hits",
		title: "title",
		url:   "url",
	}
	return &e
}

//NewStackOverflowAPI requires a config to for the url.
func NewStackOverflowAPI(config *viper.Viper) *ExternalAPI {
	var e ExternalAPI

	e.name = "Stack Overflow"
	e.url = config.GetString("StackOverflow.url")
	e.queryKey = "query"
	e.params = []URLParam{
		URLParam{"order", "desc"},
		URLParam{"sort", "relevance"},
		URLParam{"accepted", "True"},
		URLParam{"site", "stackoverflow"},
		URLParam{"pagesize", "15"},
	}
	e.parsing = ParsingParams{
		items: "items",
		title: "title",
		url:   "link",
	}
	return &e
}
