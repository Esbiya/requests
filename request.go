package requests

import (
	"errors"
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

func Get(url string, args RequestArgs) *Response {
	session := NewSession()
	return session.Get(url, args)
}

func AsyncGet(url string, args RequestArgs, ch chan *Response) {
	session := NewSession()
	go session.AsyncGet(url, args, ch)
}

func Post(url string, args RequestArgs) *Response {
	session := NewSession()
	return session.Post(url, args)
}

func AsyncPost(url string, args RequestArgs, ch chan *Response) {
	session := NewSession()
	go session.AsyncPost(url, args, ch)
}

func Head(url string, args RequestArgs) *Response {
	session := NewSession()
	return session.Head(url, args)
}

func AsyncHead(url string, args RequestArgs, ch chan *Response) {
	session := NewSession()
	go session.AsyncHead(url, args, ch)
}

func Delete(url string, args RequestArgs) *Response {
	session := NewSession()
	return session.Delete(url, args)
}

func AsyncDelete(url string, args RequestArgs, ch chan *Response) {
	session := NewSession()
	go session.AsyncDelete(url, args, ch)
}

func Options(url string, args RequestArgs) *Response {
	session := NewSession()
	return session.Options(url, args)
}

func AsyncOption(url string, args RequestArgs, ch chan *Response) {
	session := NewSession()
	go session.AsyncOptions(url, args, ch)
}

func Put(url string, args RequestArgs) *Response {
	session := NewSession()
	return session.Put(url, args)
}

func AsyncPut(url string, args RequestArgs, ch chan *Response) {
	session := NewSession()
	go session.AsyncPut(url, args, ch)
}

func Patch(url string, args RequestArgs) *Response {
	session := NewSession()
	return session.Patch(url, args)
}

func AsyncPatch(url string, args RequestArgs, ch chan *Response) {
	session := NewSession()
	go session.AsyncPatch(url, args, ch)
}
