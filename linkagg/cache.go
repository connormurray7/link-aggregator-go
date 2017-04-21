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

type LinkAggCache struct {
	redisConn redis.Conn
}

func NewLinkAggCache(config viper.Viper) *LinkAggCache {
	var c LinkAggCache
	var conn, _ = redis.Dial("tcp", config.GetString("redis.port"))
	c.redisConn = conn
	return &c
}

func (cache *LinkAggCache) Get(key string) string {
	cache.redisConn.Send("GET", key)
	val, _ := cache.redisConn.Receive()
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func (cache *LinkAggCache) Set(key string, val string) {
	cache.redisConn.Send("SET", key, val)
	cache.redisConn.Flush()
}
