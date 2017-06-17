package linkagg

import "github.com/spf13/viper"

type ExternalApi struct {
	name     string
	queryKey string
	url      string
	params   []EncodingPair
	parsing  ParsingParams
}

type EncodingPair struct {
	Key   string
	Value string
}

type ParsingParams struct {
	items string
	title string
	url   string
}

func NewGithubApi(config *viper.Viper) *ExternalApi {
	var e ExternalApi

	e.name = "Github"
	e.url = config.GetString("Github.url")
	e.queryKey = "q"
	e.params = []EncodingPair{
		EncodingPair{"sort", "stars"},
		EncodingPair{"per_page", "15"},
	}
	e.parsing = ParsingParams{
		items: "items",
		title: "name",
		url:   "html_url",
	}

	return &e
}

func NewHackerNewsApi(config *viper.Viper) *ExternalApi {
	var e ExternalApi

	e.name = "Hacker News"
	e.url = config.GetString("HackerNews.url")
	e.queryKey = "query"
	e.params = []EncodingPair{
		EncodingPair{"tags", "story"},
		EncodingPair{"hitsPerPage", "15"},
	}
	e.parsing = ParsingParams{
		items: "hits",
		title: "title",
		url:   "url",
	}

	return &e
}

func NewStackOverflowApi(config *viper.Viper) *ExternalApi {
	var e ExternalApi

	return &e
}
