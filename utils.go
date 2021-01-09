package requests

import (
	"github.com/mitchellh/mapstructure"
	"github.com/wuxiaoxiaoshen/fakeuseragent/application"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func TransferCookies(_cookies []map[string]interface{}) []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	for _, cookie := range _cookies {
		var _cookie http.Cookie
		_ = mapstructure.Decode(cookie, &_cookie)
		if _, ok := cookie["Expire"]; ok {
			_cookie.Expires = time.Unix(cookie["Expire"].(int64), 0)
		} else {
			_cookie.Expires = time.Now()
		}
		cookies = append(cookies, &_cookie)
	}
	return cookies
}

func Cookie2Map(cookie *http.Cookie) map[string]interface{} {
	var _cookie map[string]interface{}
	log.Println(cookie.MaxAge)
	_ = mapstructure.Decode(cookie, &_cookie)
	_cookie["Expires"] = (*cookie).Expires.Unix()
	return _cookie
}

func ParseHost(_url string) string {
	Url, _ := url.Parse(_url)
	return Url.Host
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func RandomNum(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

func RandomUserAgent(t interface{}) string {
	fakeUserAgent := application.NewFakeUserAgent(true, false, true)
	if t == nil {
		t = RandomNum(1, 4)
	}
	switch t {
	case 1:
		return fakeUserAgent.Chrome()
	case 2:
		return fakeUserAgent.Safari()
	case 3:
		return fakeUserAgent.IE()
	case 4:
		return fakeUserAgent.Opera()
	default:
		return fakeUserAgent.Random()
	}
}
