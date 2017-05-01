package linkagg

import (
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

//Cacher connects to a redis node and encapsulates caching.
type Cacher interface {
	Get(key string)
	Set(key string, val string)
}

//Cache holds the redis information and implements Cacher interface.
type Cache struct {
	redisConn redis.Conn
}

//NewLinkAggCache generates a new cache.
func NewLinkAggCache(config viper.Viper) *Cache {
	var c Cache
	var conn, _ = redis.Dial("tcp", config.GetString("redis.port"))
	c.redisConn = conn
	return &c
}

//Get fetches entry from Redis instance if it exists, else returns "".
func (cache *Cache) Get(key string) string {
	cache.redisConn.Send("GET", key)
	val, _ := cache.redisConn.Receive()
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

//Set saves to cache and overwrites any previous value.
func (cache *Cache) Set(key string, val string) {
	cache.redisConn.Send("SET", key, val)
	cache.redisConn.Flush()
}
