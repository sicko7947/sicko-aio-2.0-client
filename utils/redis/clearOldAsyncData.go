package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

func ClearOldAsynqData() {
	con := pool.Get()
	defer con.Close()

	if value, err := redis.Strings(con.Do("keys", "*")); err == nil {
		for _, v := range value {
			if strings.Contains(v, "asynq") {
				con.Do("DEL", v)
			}
		}
	}
}

func CheckExpireCookieInRedis() {
	con := pool.Get()
	defer con.Close()

	timestamp := fmt.Sprintf(`%d`, time.Now().Unix())
	// redis.Int64(con.Do("ZREMRANGEBYSCORE", "akamai.nike.com", "0", timestamp))
	redis.Int64(con.Do("ZREMRANGEBYSCORE", "akamai.finishline.com", "0", timestamp))
	redis.Int64(con.Do("ZREMRANGEBYSCORE", "akamai.macys.com", "0", timestamp))
	redis.Int64(con.Do("ZREMRANGEBYSCORE", "akamai.luisaviaroma.com", "0", timestamp))
}
