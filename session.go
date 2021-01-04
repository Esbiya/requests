package requests

import (
	"crypto/tls"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

type (
	CookieStore struct {
		sync.RWMutex
		v []*http.Cookie
	}
	Session struct {
		sync.Mutex
		Client                 *http.Client
		request                *http.Request
		cookieStore            *CookieStore
		beforeRequestHookFuncs []BeforeRequestHookFunc
		afterResponseHookFuncs []AfterResponseHookFunc
		option                 []ModifySessionOption
	}
)

var (
	disableRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	defaultCheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
)

func (c *CookieStore) Append(c1 *http.Cookie) {
	for i, c2 := range c.v {
		if c1.Name == c2.Name {
			c.v = append(c.v[:i], c.v[i+1:]...)
		}
	}
	c.v = append(c.v, c1)
}

type SessionArgs struct {
	proxy              string
	timeout            time.Duration
	skipVerifyTLS      bool
	chunked            bool
	allowRedirects     bool
	disableKeepAlive   bool
	disableCompression bool
}

type ModifySessionOption func(session *SessionArgs)

func Proxy(_proxy string) ModifySessionOption {
	return func(r *SessionArgs) {
		r.proxy = _proxy
	}
}

func Timeout(_timeout time.Duration) ModifySessionOption {
	return func(r *SessionArgs) {
		r.timeout = _timeout
	}
}

func SkipVerifyTLS(_skipVerifyTLS bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.skipVerifyTLS = _skipVerifyTLS
	}
}

func Chunked(_chunked bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.chunked = _chunked
	}
}

func DisableKeepAlive(_disableKeepAlive bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.disableKeepAlive = _disableKeepAlive
	}
}

func DisableCompression(_disableCompression bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.disableCompression = _disableCompression
	}
}

func AllowRedirects(_allowRedirects bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.allowRedirects = _allowRedirects
	}
}

func NewSession(opts ...ModifySessionOption) *Session {
	opt := SessionArgs{
		proxy:              "",
		timeout:            30 * time.Second,
		skipVerifyTLS:      false,
		disableKeepAlive:   false,
		disableCompression: false,
		allowRedirects:     true,
	}

	for _, f := range opts {
		f(&opt)
	}

	tranSport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   opt.timeout,
			KeepAlive: opt.timeout,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opt.skipVerifyTLS,
		},
		DisableKeepAlives:  opt.disableKeepAlive,
		DisableCompression: opt.disableCompression,
	}
	if opt.proxy != "" {
		Url, _ := url.Parse(opt.proxy)
		proxyUrl := http.ProxyURL(Url)
		tranSport.Proxy = proxyUrl
	}

	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	client.Transport = tranSport

	if opt.allowRedirects {
		client.CheckRedirect = defaultCheckRedirect
	} else {
		client.CheckRedirect = disableRedirect
	}

	return &Session{
		Client: client,
		cookieStore: &CookieStore{
			v: make([]*http.Cookie, 0),
		},
		option: opts,
	}
}

func (s *Session) InitCookieStore(_url string, cookies []*http.Cookie) {
	s.cookieStore.v = cookies
	Url, _ := url.Parse(_url)
	s.Client.Jar.SetCookies(Url, cookies)
}

func Cookie2Map(cookie *http.Cookie) map[string]interface{} {
	var _cookie map[string]interface{}
	_ = mapstructure.Decode(cookie, &_cookie)
	_cookie["Expires"] = map[string]int64{
		"timestamp": (*cookie).Expires.Unix(),
	}
	return _cookie
}

func (s *Session) CookieStore() []map[string]interface{} {
	cookies := make([]map[string]interface{}, 0)
	for _, cookie := range s.cookieStore.v {
		cookies = append(cookies, Cookie2Map(cookie))
	}
	return cookies
}

func (s *Session) CookiesMap() map[string]interface{} {
	cookies := map[string]interface{}{}
	for _, cookie := range s.cookieStore.v {
		cookies[(*cookie).Name] = (*cookie).Value
	}
	return cookies
}

func (s *Session) Request(method string, urlStr string, option Option) *Response {
	s.Lock()
	defer s.Unlock()

	method = strings.ToUpper(method)
	switch method {
	case HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH:
		urlStrParsed, err := url.Parse(urlStr)
		if err != nil {
			return &Response{
				Err: err,
			}
		}
		urlStrParsed.RawQuery = urlStrParsed.Query().Encode()

		s.request, err = http.NewRequest(method, urlStrParsed.String(), nil)
		if err != nil {
			return &Response{
				Err: err,
			}
		}
		s.request.Header.Set("User-Agent", userAgent)
		// 是否保持 keep-alive, true 表示请求完毕后关闭 tcp 连接, 不再复用
		//s.request.Close = true

		if s.Client == nil {
			s.Client = &http.Client{}
			jar, _ := cookiejar.New(nil)
			s.Client.Jar = jar
			s.Client.Transport = &http.Transport{}
		}

		if option != nil {
			err = option.setRequestOpt(s.request)
			if err != nil {
				return &Response{
					Err: err,
				}
			}

			err = option.setClientOpt(s.Client)
			if err != nil {
				return &Response{
					Err: err,
				}
			}
		}

		for _, fn := range s.beforeRequestHookFuncs {
			err = fn(s.request)
			if err != nil {
				break
			}
		}

	default:
		return &Response{
			Err: ErrInvalidMethod,
		}
	}

	r, err := s.Client.Do(s.request)
	if err != nil {
		return &Response{
			Err: err,
		}
	}

	for _, fn := range s.afterResponseHookFuncs {
		err = fn(r)
		if err != nil {
			break
		}
	}

	for _, cookie := range r.Cookies() {
		s.cookieStore.Lock()
		s.cookieStore.Append(cookie)
		s.cookieStore.Unlock()
	}

	return NewResponse(r)
}

func (s *Session) AsyncRequest(method string, urlStr string, option Option, ch chan *Response) {
	response := s.Request(method, urlStr, option)
	ch <- response
}

func (s *Session) GetRequest() *http.Request {
	return s.request
}

type (
	BeforeRequestHookFunc func(*http.Request) error
	AfterResponseHookFunc func(*http.Response) error
)

func (s *Session) RegisterBeforeReqHook(fn BeforeRequestHookFunc) error {
	if s.beforeRequestHookFuncs == nil {
		s.beforeRequestHookFuncs = make([]BeforeRequestHookFunc, 0, 8)
	}
	if len(s.beforeRequestHookFuncs) > 7 {
		return ErrHookFuncMaxLimit
	}
	s.beforeRequestHookFuncs = append(s.beforeRequestHookFuncs, fn)
	return nil
}

func (s *Session) UnregisterBeforeReqHook(index int) error {
	if index >= len(s.beforeRequestHookFuncs) {
		return ErrIndexOutofBound
	}
	s.beforeRequestHookFuncs = append(s.beforeRequestHookFuncs[:index], s.beforeRequestHookFuncs[index+1:]...)
	return nil
}

func (s *Session) ResetBeforeReqHook() {
	s.beforeRequestHookFuncs = []BeforeRequestHookFunc{}
}

func (s *Session) RegisterAfterRespHook(fn AfterResponseHookFunc) error {
	if s.afterResponseHookFuncs == nil {
		s.afterResponseHookFuncs = make([]AfterResponseHookFunc, 0, 8)
	}
	if len(s.afterResponseHookFuncs) > 7 {
		return ErrHookFuncMaxLimit
	}
	s.afterResponseHookFuncs = append(s.afterResponseHookFuncs, fn)
	return nil
}

func (s *Session) Copy(_url string) *Session {
	opt := s.option
	session := NewSession(opt...)
	session.cookieStore = s.cookieStore
	session.InitCookieStore(_url, s.cookieStore.v)
	return session
}

func (s *Session) UnregisterAfterRespHook(index int) error {
	if index >= len(s.afterResponseHookFuncs) {
		return ErrIndexOutofBound
	}
	s.afterResponseHookFuncs = append(s.afterResponseHookFuncs[:index], s.afterResponseHookFuncs[index+1:]...)
	return nil
}

func (s *Session) ResetAfterRespHook() {
	s.afterResponseHookFuncs = []AfterResponseHookFunc{}
}

func (s *Session) Get(url string, option Option) *Response {
	return s.Request("get", url, option)
}

func (s *Session) AsyncGet(url string, option Option, ch chan *Response) {
	go s.AsyncRequest("get", url, option, ch)
}

func (s *Session) Post(url string, option Option) *Response {
	return s.Request("post", url, option)
}

func (s *Session) AsyncPost(url string, option Option, ch chan *Response) {
	go s.AsyncRequest("post", url, option, ch)
}

func (s *Session) Head(url string, option Option) *Response {
	return s.Request("head", url, option)
}

func (s *Session) AsyncHead(url string, option Option, ch chan *Response) {
	go s.AsyncRequest("head", url, option, ch)
}

func (s *Session) Delete(url string, option Option) *Response {
	return s.Request("delete", url, option)
}

func (s *Session) AsyncDelete(url string, option Option, ch chan *Response) {
	go s.AsyncRequest("delete", url, option, ch)
}

func (s *Session) Options(url string, option Option) *Response {
	return s.Request("options", url, option)
}

func (s *Session) AsyncOptions(url string, option Option, ch chan *Response) {
	go s.AsyncRequest("options", url, option, ch)
}

func (s *Session) Put(url string, option Option) *Response {
	return s.Request("put", url, option)
}

func (s *Session) AsyncPut(url string, option Option, ch chan *Response) {
	go s.AsyncRequest("put", url, option, ch)
}

func (s *Session) Patch(url string, option Option) *Response {
	return s.Request("patch", url, option)
}

func (s *Session) AsyncPatch(url string, option Option, ch chan *Response) {
	go s.AsyncRequest("patch", url, option, ch)
}
