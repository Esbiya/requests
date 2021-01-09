package requests

import (
	"crypto/tls"
	"github.com/pkg/errors"
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
		request                *http.Request
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

type SessionArgs struct {
	url                *url.URL
	cookies            []*http.Cookie
	proxy              string
	timeout            time.Duration
	skipVerifyTLS      bool
	chunked            bool
	allowRedirects     bool
	disableKeepAlive   bool
	disableCompression bool
}

type ModifySessionOption func(session *SessionArgs)

func Url(_url string) ModifySessionOption {
	return func(r *SessionArgs) {
		r.url, _ = url.Parse(_url)
	}
}

func Proxy(_proxy string) ModifySessionOption {
	return func(r *SessionArgs) {
		r.proxy = _proxy
	}
}

func Cookies(_cookies []map[string]interface{}) ModifySessionOption {
	return func(r *SessionArgs) {
		r.cookies = TransferCookies(_cookies)
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

	cookieJar := NewCookieJar()

	client := &http.Client{}
	client.Jar = cookieJar
	client.Transport = tranSport

	if opt.allowRedirects {
		client.CheckRedirect = defaultCheckRedirect
	} else {
		client.CheckRedirect = disableRedirect
	}

	session := &Session{
		Url:       opt.url,
		Client:    client,
		option:    opts,
		CookieJar: cookieJar,
	}

	if session.Url != nil && opt.cookies != nil {
		session.CookieJar.SetCookies(session.Url, opt.cookies)
	}
	return session
}

func (s *Session) SetUrl(_url string) *Session {
	s.Url, _ = url.Parse(_url)
	return s
}

func (s *Session) SetCookies(_url string, cookies []*http.Cookie) *Session {
	Url, err := url.Parse(_url)
	if err != nil {
		panic("set cookies error: " + err.Error())
	}
	s.CookieJar.SetCookies(Url, cookies)
	return s
}

func (s *Session) Cookies(_url string) []*http.Cookie {
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
	return s
}

func (s *Session) SetSkipVerifyTLS(ssl bool) *Session {
	s.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = ssl
	return s
}

func (s *Session) SetProxy(proxy string) *Session {
	proxyUrl, _ := url.Parse(proxy)
	s.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
	return s
}

func (s *Session) SetDisableKeepAlive(disable bool) *Session {
	s.Client.Transport.(*http.Transport).DisableKeepAlives = disable
	return s
}

func (s *Session) SetDisableCompression(disable bool) *Session {
	s.Client.Transport.(*http.Transport).DisableCompression = disable
	return s
}

func (s *Session) SetAllowRedirect(y bool) *Session {
	if y {
		s.Client.CheckRedirect = defaultCheckRedirect
	} else {
		s.Client.CheckRedirect = disableRedirect
	}
	return s
}

func (s *Session) Save(path string, _url string) error {
	return s.CookieJar.Save(path, _url)
}

func (s *Session) Load(path string, _url string) error {
	return s.CookieJar.Load(path, _url)
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
			s.CookieJar = NewCookieJar()
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

func (s *Session) Copy() *Session {
	opt := s.option
	session := NewSession(opt...)
	session.CookieJar = s.CookieJar
	return session
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
