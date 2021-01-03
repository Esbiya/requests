package requests

import (
	"errors"
)

const (
	userAgent = "requests 1.0"
)

var (
	// ErrInvalidMethod will be throwed when method not in
	// [HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH, CONNECT, TRACE]
	ErrInvalidMethod = errors.New("nic: Method is invalid")

	// ErrFileInfo will be throwed when fileinfo is invalid
	ErrFileInfo = errors.New("nic: Invalid file information")

	// ErrParamConflict will be throwed when options params conflict
	// e.g. files + data
	//      json + data
	//      ...
	ErrParamConflict = errors.New("nic: RequestArgs param conflict")

	// ErrUnrecognizedEncoding will be throwed while changing response encoding
	// if encoding is not recognized
	ErrUnrecognizedEncoding = errors.New("nic: Unrecognized encoding")

	// ErrNotJsonResponse will be throwed when response not a json
	// but invoke Json() method
	ErrNotJsonResponse = errors.New("nic: Not a Json response")

	// ErrHookFuncMaxLimit will be throwed when the number of hook functions
	// more than MaxLimit = 8
	ErrHookFuncMaxLimit = errors.New("nic: The number of hook functions must be less than 8")

	// ErrIndexOutofBound means the index out of bound
	ErrIndexOutofBound = errors.New("nic: Index out of bound")
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
