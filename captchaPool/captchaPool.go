package captchaPool

import (
	"time"

	"github.com/gogf/gf/container/gpool"
)

var (
	supplystoreCaptchaPool *gpool.Pool
	adidasCaptchaPool      *gpool.Pool
)

func init() {
	supplystoreCaptchaPool = gpool.New(time.Duration(2*time.Minute), nil)
	adidasCaptchaPool = gpool.New(time.Duration(2*time.Minute), nil)
}

func GetSupplyStoreCaptchaToken() (token string, err error) {
	var value interface{}
	if value, err = supplystoreCaptchaPool.Get(); err != nil {
		return "", err
	}
	return value.(string), nil
}

func GetAdidasCaptchaToken() (token string, err error) {
	var value interface{}
	if value, err = adidasCaptchaPool.Get(); err != nil {
		return "", err
	}
	return value.(string), nil
}

func PutSupplyStoreCaptchaToken(token string) {
	supplystoreCaptchaPool.Put(token)
}

func PutAdidasCaptchaToken(token string) {
	adidasCaptchaPool.Put(token)
}
