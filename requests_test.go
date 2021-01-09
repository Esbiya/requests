package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	session := NewSession()
	resp := session.Get("https://www.baidu.com", RequestArgs{})
	log.Println(resp.Cookies())
	c1, _ := session.CookieJar.Array("")
	x, _ := json.MarshalIndent(c1, "", "    ")
	log.Println(string(x))

	log.Println(session.CookieJar.String(""))

	session.SetCookies("https://www.baidu.com", []*http.Cookie{
		{
			Name:  "BIDUPSID",
			Value: "BIDUPSID",
		},
	})
	log.Println(session.CookieJar.String(""))

	session.SetProxy("http://127.0.0.1:8888")
	session.Get("https://www.baidu.com", RequestArgs{})
	log.Println(session.CookieJar.String("https://www.baidu.com"))
}
