package requests

import (
	"bytes"
	"crypto/tls"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type (
	Session struct {
		sync.Mutex
		Url                    *url.URL
		Client                 *http.Client
		CookieJar              *CookieJar
		request                *Request
		beforeRequestHookFuncs []BeforeRequestHookFunc
		afterResponseHookFuncs []AfterResponseHookFunc
		option                 []ModifySessionOption
		args                   Arguments
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

type Arguments struct {
	Url                string
	Cookies            []*http.Cookie
	Proxy              string
	Timeout            time.Duration
	SkipVerifyTLS      bool
	Chunked            bool
	AllowRedirects     bool
	DisableKeepAlive   bool
	DisableCompression bool
}

func (h *Arguments) setClientArgs(client *http.Client) error {
	if !h.AllowRedirects {
		client.CheckRedirect = disableRedirect
	} else {
		client.CheckRedirect = defaultCheckRedirect
	}

	client.Timeout = h.Timeout

	transport := client.Transport.(*http.Transport)
	transport.DisableKeepAlives = h.DisableKeepAlive
	transport.DisableCompression = h.DisableCompression

	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}

	transport.TLSClientConfig.InsecureSkipVerify = h.SkipVerifyTLS

	if h.Proxy != "" {
		proxyUrl, err := url.Parse(h.Proxy)
		if err != nil {
			return err
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	return nil
}

type ModifySessionOption func(session *Arguments)

func SetUrl(_url string) ModifySessionOption {
	return func(r *Arguments) {
		r.Url = _url
	}
}

func SetProxy(_proxy string) ModifySessionOption {
	return func(r *Arguments) {
		r.Proxy = _proxy
	}
}

func SetCookies(_cookies []*http.Cookie) ModifySessionOption {
	return func(r *Arguments) {
		r.Cookies = _cookies
	}
}

func SetTimeout(_timeout time.Duration) ModifySessionOption {
	return func(r *Arguments) {
		r.Timeout = _timeout
	}
}

func SetSkipVerifyTLS(_skipVerifyTLS bool) ModifySessionOption {
	return func(r *Arguments) {
		r.SkipVerifyTLS = _skipVerifyTLS
	}
}

func SetDisableKeepAlive(_disableKeepAlive bool) ModifySessionOption {
	return func(r *Arguments) {
		r.DisableKeepAlive = _disableKeepAlive
	}
}

func SetDisableCompression(_disableCompression bool) ModifySessionOption {
	return func(r *Arguments) {
		r.DisableCompression = _disableCompression
	}
}

func NewSession(opts ...ModifySessionOption) *Session {
	args := Arguments{
		Proxy:              "",
		Timeout:            30 * time.Second,
		SkipVerifyTLS:      false,
		DisableKeepAlive:   false,
		DisableCompression: false,
		AllowRedirects:     true,
	}

	for _, f := range opts {
		f(&args)
	}

	tranSport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   args.Timeout,
			KeepAlive: args.Timeout,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: args.SkipVerifyTLS,
		},
		DisableKeepAlives:     args.DisableKeepAlive,
		DisableCompression:    args.DisableCompression,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if args.Proxy != "" {
		Url, _ := url.Parse(args.Proxy)
		proxyUrl := http.ProxyURL(Url)
		tranSport.Proxy = proxyUrl
	}

	client := &http.Client{}
	client.Transport = tranSport

	if args.AllowRedirects {
		client.CheckRedirect = defaultCheckRedirect
	} else {
		client.CheckRedirect = disableRedirect
	}

	Url, _ := url.Parse(args.Url)
	session := &Session{
		Url:       Url,
		option:    opts,
		CookieJar: NewCookieJar(),
		args:      args,
	}
	client.Jar = session.CookieJar
	session.Client = client

	if session.Url != nil && args.Cookies != nil {
		session.CookieJar.SetCookies(session.Url, args.Cookies)
	}
	return session
}

func (s *Session) SetUrl(_url string) *Session {
	s.Url, _ = url.Parse(_url)
	s.args.Url = _url
	return s
}

func (s *Session) GetUrl() string {
	return s.args.Url
}

func (s *Session) SetCookies(_url string, cookies []*http.Cookie) *Session {
	Url, err := url.Parse(_url)
	if err != nil {
		panic("set Cookies error: " + err.Error())
	}
	s.CookieJar.SetCookies(Url, cookies)
	return s
}

func (s *Session) GetCookies(_url string) []*http.Cookie {
	var Url *url.URL
	if _url == "" {
		Url = s.Url
	} else {
		Url, _ = url.Parse(_url)
	}
	if Url == nil {
		return []*http.Cookie{}
	}
	return s.CookieJar.Cookies(Url)
}

func (s *Session) SetTimeout(timeout time.Duration) *Session {
	s.Client.Timeout = timeout
	s.args.Timeout = timeout
	return s
}

func (s *Session) GetTimeout() time.Duration {
	return s.args.Timeout
}

func (s *Session) SetSkipVerifyTLS(skip bool) *Session {
	s.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = skip
	s.args.SkipVerifyTLS = skip
	return s
}

func (s *Session) GetSkipVerifyTLS() bool {
	return s.args.SkipVerifyTLS
}

func (s *Session) SetProxy(proxy string) *Session {
	proxyUrl, _ := url.Parse(proxy)
	s.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
	s.args.Proxy = proxy
	return s
}

func (s *Session) GetProxy() string {
	return s.args.Proxy
}

func (s *Session) SetDisableKeepAlive(disable bool) *Session {
	s.Client.Transport.(*http.Transport).DisableKeepAlives = disable
	return s
}

func (s *Session) GetDisableKeepAlive() bool {
	return s.args.DisableKeepAlive
}

func (s *Session) SetDisableCompression(disable bool) *Session {
	s.Client.Transport.(*http.Transport).DisableCompression = disable
	return s
}

func (s *Session) GetDisableCompression() bool {
	return s.args.DisableCompression
}

func (s *Session) SetAllowRedirects(y bool) *Session {
	if y {
		s.Client.CheckRedirect = defaultCheckRedirect
	} else {
		s.Client.CheckRedirect = disableRedirect
	}
	return s
}

func (s *Session) GetAllowRedirects() bool {
	return s.args.AllowRedirects
}

func (s *Session) Save(path string, _url string) error {
	return s.CookieJar.Save(path, _url)
}

func (s *Session) Load(path string, _url string) error {
	return s.CookieJar.Load(path, _url)
}

func (s *Session) AsyncDo(method string, urlStr string, ch chan *Response, args ...interface{}) {
	resp := s.Do(method, urlStr, args...)
	ch <- resp
}

func (s *Session) Do(method string, urlStr string, args ...interface{}) *Response {
	s.Lock()
	defer s.Unlock()

	method = strings.ToUpper(method)
	switch method {
	case HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH:

		urlStrParsed, err := url.Parse(urlStr)
		if err != nil {
			return &Response{
				err: err,
			}
		}
		urlStrParsed.RawQuery = urlStrParsed.Query().Encode()

		s.request = &Request{Files: make([]*File, 0)}
		s.request.Request, err = http.NewRequest(method, urlStrParsed.String(), nil)
		if err != nil {
			return &Response{
				err: err,
			}
		}
		s.request.Header.Set("User-Agent", RandomUserAgent(nil))

		if s.Client == nil {
			s.CookieJar = NewCookieJar()
			s.Client = &http.Client{}
			s.Client.Jar = s.CookieJar
			s.Client.Transport = &http.Transport{}
		}

		for _, arg := range args {
			switch _arg := arg.(type) {
			case Headers:
				s.request.Headers = _arg
			case SimpleCookie:
				s.request.Cookies = _arg
			case Auth:
				s.request.Auth = _arg
			case Params:
				s.request.Params = _arg
			case Form:
				s.request.Form = _arg
			case Payload:
				s.request.Payload = _arg
			case *File:
				s.request.Files = append(s.request.Files, _arg)
			case []*File:
				s.request.Files = append(s.request.Files, _arg...)
			case string:
				body := strings.NewReader(_arg)
				s.request.Request.Body = ioutil.NopCloser(body)
				if !s.request.Chunked {
					s.request.Request.ContentLength = int64(len(_arg))
				}
			case []byte:
				body := bytes.NewReader(_arg)
				s.request.Request.Body = ioutil.NopCloser(body)
				if !s.request.Chunked {
					s.request.Request.ContentLength = int64(len(_arg))
				}
			case *http.Cookie:
				s.request.Request.AddCookie(_arg)
			case []*http.Cookie:
				for _, cookie := range _arg {
					s.request.Request.AddCookie(cookie)
				}
			case Arguments:
				err := _arg.setClientArgs(s.Client)
				if err != nil {
					return &Response{
						err: err,
					}
				}
			}
		}

		for _, fn := range s.beforeRequestHookFuncs {
			err = fn(s.request)
			if err != nil {
				break
			}
		}

		err = s.request.setRequestArgs()
		if err != nil {
			return &Response{
				err: err,
			}
		}

	default:
		return &Response{
			err: ErrInvalidMethod,
		}
	}

	before := time.Now()
	r, err := s.Client.Do(s.request.Request)
	after := time.Now()
	cost := after.Sub(before)

	if err != nil {
		return &Response{
			cost: cost,
			err:  err,
		}
	}

	resp := NewResponse(r, cost)

	for _, fn := range s.afterResponseHookFuncs {
		err = fn(resp)
		if err != nil {
			break
		}
	}
	return resp
}

func (s *Session) GetRequest() *Request {
	return s.request
}

func (s *Session) Copy() *Session {
	opt := s.option
	session := NewSession(opt...)
	session.CookieJar = s.CookieJar
	session.Client.Jar = s.CookieJar
	return session
}

type (
	BeforeRequestHookFunc func(*Request) error
	AfterResponseHookFunc func(*Response) error
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
		return ErrIndexOutOfBound
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

func (s *Session) UnregisterAfterRespHook(index int) error {
	if index >= len(s.afterResponseHookFuncs) {
		return ErrIndexOutOfBound
	}
	s.afterResponseHookFuncs = append(s.afterResponseHookFuncs[:index], s.afterResponseHookFuncs[index+1:]...)
	return nil
}

func (s *Session) ResetAfterRespHook() {
	s.afterResponseHookFuncs = []AfterResponseHookFunc{}
}

func (s *Session) Get(url string, args ...interface{}) *Response {
	return s.Do("get", url, args...)
}

func (s *Session) AsyncGet(url string, args ...interface{}) *AsyncResponse {
	c := make(chan *Response, 1)
	go s.AsyncDo("get", url, c, args...)
	return &AsyncResponse{wg: &sync.WaitGroup{}, c: c}
}

func (s *Session) Post(url string, args ...interface{}) *Response {
	return s.Do("post", url, args...)
}

func (s *Session) AsyncPost(url string, args ...interface{}) *AsyncResponse {
	c := make(chan *Response, 1)
	go s.AsyncDo("post", url, c, args...)
	return &AsyncResponse{wg: &sync.WaitGroup{}, c: c}
}

func (s *Session) Head(url string, args ...interface{}) *Response {
	return s.Do("head", url, args...)
}

func (s *Session) AsyncHead(url string, args ...interface{}) *AsyncResponse {
	c := make(chan *Response, 1)
	go s.AsyncDo("head", url, c, args...)
	return &AsyncResponse{wg: &sync.WaitGroup{}, c: c}
}

func (s *Session) Delete(url string, args ...interface{}) *Response {
	return s.Do("delete", url, args...)
}

func (s *Session) AsyncDelete(url string, args ...interface{}) *AsyncResponse {
	c := make(chan *Response, 1)
	go s.AsyncDo("delete", url, c, args...)
	return &AsyncResponse{wg: &sync.WaitGroup{}, c: c}
}

func (s *Session) Options(url string, args ...interface{}) *Response {
	return s.Do("options", url, args...)
}

func (s *Session) AsyncOptions(url string, args ...interface{}) *AsyncResponse {
	c := make(chan *Response, 1)
	go s.AsyncDo("option", url, c, args...)
	return &AsyncResponse{wg: &sync.WaitGroup{}, c: c}
}

func (s *Session) Put(url string, args ...interface{}) *Response {
	return s.Do("put", url, args...)
}

func (s *Session) AsyncPut(url string, args ...interface{}) *AsyncResponse {
	c := make(chan *Response, 1)
	go s.AsyncDo("put", url, c, args...)
	return &AsyncResponse{wg: &sync.WaitGroup{}, c: c}
}

func (s *Session) Patch(url string, args ...interface{}) *Response {
	return s.Do("patch", url, args...)
}

func (s *Session) AsyncPatch(url string, args ...interface{}) *AsyncResponse {
	c := make(chan *Response, 1)
	go s.AsyncDo("patch", url, c, args...)
	return &AsyncResponse{wg: &sync.WaitGroup{}, c: c}
}
