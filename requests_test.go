package requests

import (
	"encoding/json"
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	session := NewSession()
	resp := session.Get("https://www.baidu.com", RequestArgs{})
	log.Println(resp.Cookies())
	x, _ := json.MarshalIndent(session.CookieJar.Map(), "", "    ")
	log.Println(string(x))
}
