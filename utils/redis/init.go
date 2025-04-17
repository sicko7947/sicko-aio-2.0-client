package redis

import (
	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

func init() {
	pool = NewPool()

	// go func() {
	// 	t := time.NewTicker(20 * time.Second)
	// 	for {
	// 		CheckExpireCookieInRedis()
	// 		<-t.C
	// 	}
	// }()
}

// NewPool : Create Redis Pool recycle connections
func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
