package requests

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	session := NewSession().SetProxy("http://127.0.0.1:8888").SetSkipVerifyTLS(false)

	url := "https://www.baidu.com"
	session.Get("https://www.baidu.com", RequestArgs{})
	_ = session.CookieJar.Save("./cookies.json", url)

	session1 := NewSession().SetProxy("http://127.0.0.1:8888")
	log.Println(session1.CookieJar == session1.Client.Jar)
	_ = session1.Load("./cookies.json", url)
	session1.Get("https://www.baidu.com", RequestArgs{})
}
