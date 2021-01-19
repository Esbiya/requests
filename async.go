package requests

import (
	"errors"
	"sync"
	"time"
)

var (
	asyncResponses = make([]*AsyncResponse, 0)
)

type AsyncResponse struct {
	wg  *sync.WaitGroup
	err error
	c   chan *Response
}

func (a *AsyncResponse) Then(f func(r *Response)) *AsyncResponse {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		select {
		case resp := <-a.c:
			a.err = resp.Error()
			f(resp)
		case <-time.After(30 * time.Second):
			a.err = errors.New("async requests timeout")
		}
	}()
	asyncResponses = append(asyncResponses, a)
	return a
}

func (a *AsyncResponse) Catch(f func(e error)) *AsyncResponse {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		f(a.err)
	}()
	return a
}

func (a *AsyncResponse) Wait() {
	a.wg.Wait()
}

func AsyncWait() {
	for _, s := range asyncResponses {
		s.Wait()
	}
}

func AsyncGet(url string, args ...interface{}) *AsyncResponse {
	session := NewSession()
	return session.AsyncGet(url, args...)
}

func AsyncPost(url string, args ...interface{}) *AsyncResponse {
	session := NewSession()
	return session.AsyncPost(url, args...)
}

func AsyncHead(url string, args ...interface{}) *AsyncResponse {
	session := NewSession()
	return session.AsyncHead(url, args...)
}

func AsyncDelete(url string, args ...interface{}) *AsyncResponse {
	session := NewSession()
	return session.AsyncDelete(url, args...)
}

func AsyncOptions(url string, args ...interface{}) *AsyncResponse {
	session := NewSession()
	return session.AsyncOptions(url, args...)
}

func AsyncPut(url string, args ...interface{}) *AsyncResponse {
	session := NewSession()
	return session.AsyncPut(url, args...)
}

func AsyncPatch(url string, args ...interface{}) *AsyncResponse {
	session := NewSession()
	return session.AsyncPatch(url, args...)
}
