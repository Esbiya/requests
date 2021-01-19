package example

import (
	"github.com/Esbiya/requests"
	"log"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	params := requests.Params{"1": "2"}
	resp := requests.Get("https://www.baidu.com", params, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true})
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
	log.Println(resp.Cost().String())
}

func TestSessionGet(t *testing.T) {
	session := requests.NewSession()
	resp := session.Get("https://www.baidu.com")
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
	log.Println(resp.Cost().String())
}

func TestAsyncGet(t *testing.T) {
	start := time.Now()
	for i := 0; i < 10; i++ {
		headers := requests.Headers{
			"Connection": "keep-alive",
		}
		requests.AsyncGet("https://www.baidu.com", headers).Then(func(r *requests.Response) {
			if r.Error() != nil {
				log.Fatal(r.Error())
			}
			log.Println("function cost => " + r.Cost().String())
		})
	}
	log.Println("execute first")
	requests.AsyncWait() // 必须等待, 否则异步请求不会执行
	end := time.Now()
	log.Println("all cost => " + end.Sub(start).String())
}

func TestSessionAsyncGet(t *testing.T) {
	start := time.Now()
	session := requests.NewSession()
	for i := 0; i < 10; i++ {
		c := make(chan *requests.Response, 1)
		headers := requests.Headers{
			"Connection": "keep-alive",
		}
		session.AsyncGet("https://www.baidu.com", c, headers).Then(func(r *requests.Response) {
			if r.Error() != nil {
				log.Fatal(r.Error())
			}
			log.Println("function cost => " + r.Cost().String())
		})
	}
	log.Println("execute first")
	requests.AsyncWait() // 必须等待, 否则异步请求不会执行
	end := time.Now()
	log.Println("all cost => " + end.Sub(start).String())
}
