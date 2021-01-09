package requests

import (
	"github.com/mitchellh/mapstructure"
	"net/http"
	"net/url"
	"os"
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
