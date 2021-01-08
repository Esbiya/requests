package requests

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	session := NewSession()
	resp := session.Get("https://www.baidu.com", RequestArgs{}).Text
	log.Println(resp)
	log.Println(session.Cookies("https://www.baidu.com"))
}
