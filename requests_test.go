package requests

import (
	"encoding/json"
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	session := NewSession()
	resp := session.Get("https://www.baidu.com", RequestArgs{}).Text
	log.Println(resp)
	x, _ := json.MarshalIndent(session.CookieJar.Map(), "", "    ")
	log.Println(string(x))
}
