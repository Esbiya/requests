package requests

import (
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func TransferCookies(_cookies []map[string]interface{}) []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	for _, cookie := range _cookies {
		var _cookie http.Cookie
		_ = mapstructure.Decode(cookie, &_cookie)
		_cookie.Expires = time.Unix(cookie["Expire"].(int64), 0)
		cookies = append(cookies, &_cookie)
	}
	return cookies
}
