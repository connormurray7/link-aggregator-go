package linkagg

//Requester makes a request to the outside APIs
type Requester interface {
	Request(req string) string
}

type LinkAgg struct {
	cache Cache
}

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

}
