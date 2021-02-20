package example

import (
	"encoding/json"
	"github.com/Esbiya/requests"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	params := requests.Params{"test": "123"}
	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		for key, value := range params {
			if v := query.Get(key); value != v {
				t.Errorf("query param %s = %s; want = %s", key, v, value)
			}
		}
		x, _ := json.Marshal(map[string]interface{}{
			"xxx": map[string]interface{}{
				"yyy": 11111,
			},
		})
		_, _ = w.Write(x)
	}
	server := httptest.NewServer(http.HandlerFunc(queryHandler))

	resp := requests.Get(server.URL, params)
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	result, _ := resp.JSON()
	log.Println(result.Get("xxx.yyy"))
	log.Println(resp.Cost().String())
}

func TestSessionGet(t *testing.T) {
	params := requests.Params{"test": "123"}
	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		for key, value := range params {
			if v := query.Get(key); value != v {
				t.Errorf("query param %s = %s; want = %s", key, v, value)
			}
		}
		_, _ = w.Write([]byte(query.Encode()))
	}
	server := httptest.NewServer(http.HandlerFunc(queryHandler))

	session := requests.NewSession()
	resp := session.Get(server.URL, params)
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
	log.Println(resp.Cost().String())
}

func TestAsyncGet(t *testing.T) {
	params := requests.Params{"test": "123"}
	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		for key, value := range params {
			if v := query.Get(key); value != v {
				t.Errorf("query param %s = %s; want = %s", key, v, value)
			}
		}
		_, _ = w.Write([]byte(query.Encode()))
	}
	server := httptest.NewServer(http.HandlerFunc(queryHandler))

	start := time.Now()
	for i := 0; i < 10; i++ {
		headers := requests.Headers{
			"Connection": "keep-alive",
		}
		requests.AsyncGet(server.URL, headers, params).Then(func(r *requests.Response) {
			if r.Error() != nil {
				log.Fatal(r.Error())
			}
			log.Println(r.Text)
			log.Println("function cost => " + r.Cost().String())
		})
	}
	log.Println("execute first")
	requests.AsyncWait() // 必须等待, 否则异步请求不会执行
	end := time.Now()
	log.Println("all cost => " + end.Sub(start).String())
}

func TestSessionAsyncGet(t *testing.T) {
	params := requests.Params{"test": "123"}
	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		for key, value := range params {
			if v := query.Get(key); value != v {
				t.Errorf("query param %s = %s; want = %s", key, v, value)
			}
		}
		_, _ = w.Write([]byte(query.Encode()))
	}
	server := httptest.NewServer(http.HandlerFunc(queryHandler))

	start := time.Now()
	session := requests.NewSession()
	for i := 0; i < 10; i++ {
		c := make(chan *requests.Response, 1)
		headers := requests.Headers{
			"Connection": "keep-alive",
		}
		session.AsyncGet(server.URL, c, headers, params).Then(func(r *requests.Response) {
			if r.Error() != nil {
				log.Fatal(r.Error())
			}
			log.Println(r.Text)
			log.Println("function cost => " + r.Cost().String())
		})
	}
	log.Println("execute first")
	requests.AsyncWait() // 必须等待, 否则异步请求不会执行
	end := time.Now()
	log.Println("all cost => " + end.Sub(start).String())
}
