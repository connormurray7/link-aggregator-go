package linkagg

import (
	"log"

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
func NewLinkAggCache(config *viper.Viper) *Cache {
	var c Cache
	log.Println("Connecting to Redis instance on local port,", config.GetString("Redis.port"))
	var conn, _ = redis.Dial("tcp", ":"+config.GetString("Redis.port"))
	c.redisConn = conn
	return &c
}

//Get fetches entry from Redis instance if it exists, else returns "".
func (cache *Cache) Get(key string) string {
	val, _ := cache.redisConn.Do("GET", key)
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
