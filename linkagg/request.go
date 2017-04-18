package linkagg

//Requester makes a request to the outside APIs
type Requester interface {
	Request(req string)
}
