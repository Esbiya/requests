package example

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Esbiya/requests"
)

func TestHeaders(t *testing.T) {
	headers := requests.Headers{
		"User-Agent": "test",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		for key, value := range headers {
			if v := r.Header.Get(key); value != v {
				t.Errorf("header %q = %s; want = %s", key, v, value)
			}
		}
		b, _ := json.Marshal(r.Header)
		_, _ = w.Write(b)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	resp := requests.Get(server.URL, headers)
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
	log.Println(resp.Cost().String())
}

func TestCookies(t *testing.T) {
	cookies := requests.SimpleCookie{
		"token": "test",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(r.Cookies())
		_, _ = w.Write(b)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	resp := requests.Get(server.URL, cookies)
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
	log.Println(resp.Cost().String())
}
