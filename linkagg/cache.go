package linkagg

//Cacher connects to a redis node and encapsulates caching.
type Cacher interface {
	Get(key string)
}
