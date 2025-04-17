package redis

import (
	"encoding/json"
	"errors"
	"strings"

	http "github.com/zMrKrabz/fhttp"

	"github.com/garyburd/redigo/redis"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/models"
)

// GetCookie2FromRedis : Get Single bm_sz Cookie From Redis
func GetCookie2FromRedis(domain string) (string, map[string]*http.Cookie, *models.Error) {
	con := pool.Get()
	defer con.Close()

	cookieList, _ := redis.Strings(con.Do("ZPOPMIN", domain))
	if len(cookieList) > 0 {
		data, err := sickocommon.LZDecompress(strings.ReplaceAll(cookieList[0], `"`, ``), "")
		if err != nil {
			return "", nil, &models.Error{Error: err}
		}

		result := gjson.Parse(data)
		useragent := result.Get("useragent").String()
		cookieMap := make(map[string]*http.Cookie)

		json.Unmarshal([]byte(result.Get("cookies").String()), &cookieMap)

		return useragent, cookieMap, nil
	}
	return "", nil, &models.Error{
		Error:   errors.New("error gettin cookies"),
		Code:    400,
		Message: "Error Getting Cookie",
	}
}
