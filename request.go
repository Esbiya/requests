package requests

import (
	"errors"
)

const (
	userAgent = "requests 1.0"
)

var (
	ErrInvalidMethod = errors.New("requests: method is invalid")

	ErrFileInfo = errors.New("requests: invalid file information")

	ErrParamConflict = errors.New("requests: requestArgs param conflict")

	ErrUnrecognizedEncoding = errors.New("requests: unrecognized encoding")

	ErrNotJSONResponse = errors.New("requests: not a json response")

	ErrHookFuncMaxLimit = errors.New("requests: the number of hook functions must be less than 8")

	ErrIndexOutOfBound = errors.New("requests: index out of bound")
)

const (
	HEAD    = "HEAD"
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
	PUT     = "PUT"
	PATCH   = "PATCH"
)

func Get(url string, option Option) *Response {
	session := NewSession()
	return session.Get(url, option)
}

func AsyncGet(url string, option Option, ch chan *Response) {
	session := NewSession()
	go session.AsyncGet(url, option, ch)
}

func Post(url string, option Option) *Response {
	session := NewSession()
	return session.Post(url, option)
}

func AsyncPost(url string, option Option, ch chan *Response) {
	session := NewSession()
	go session.AsyncPost(url, option, ch)
}

func Head(url string, option Option) *Response {
	session := NewSession()
	return session.Head(url, option)
}

func AsyncHead(url string, option Option, ch chan *Response) {
	session := NewSession()
	go session.AsyncHead(url, option, ch)
}

func Delete(url string, option Option) *Response {
	session := NewSession()
	return session.Delete(url, option)
}

func AsyncDelete(url string, option Option, ch chan *Response) {
	session := NewSession()
	go session.AsyncDelete(url, option, ch)
}

func Options(url string, option Option) *Response {
	session := NewSession()
	return session.Options(url, option)
}

func AsyncOption(url string, option Option, ch chan *Response) {
	session := NewSession()
	go session.AsyncOptions(url, option, ch)
}

func Put(url string, option Option) *Response {
	session := NewSession()
	return session.Put(url, option)
}

func AsyncPut(url string, option Option, ch chan *Response) {
	session := NewSession()
	go session.AsyncPut(url, option, ch)
}

func Patch(url string, option Option) *Response {
	session := NewSession()
	return session.Patch(url, option)
}

func AsyncPatch(url string, option Option, ch chan *Response) {
	session := NewSession()
	go session.AsyncPatch(url, option, ch)
}
