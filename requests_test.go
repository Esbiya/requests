package requests

import (
	"log"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	url := "https://www.baidu.com/"
	resp := Get(url, RequestArgs{})
	if resp.StatusCode != http.StatusOK {
		log.Fatal("状态码异常")
	}
	log.Println(resp.Text)

	session := NewSession()
	session.CookieJar.Get("https://jd.com/").String()
}
