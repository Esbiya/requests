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
		Url                        *url.URL
		Client                     *http.Client
		CookieJar                  *CookieJar
		request                    *http.Request
		beforeRequestHookFuncs     []BeforeRequestHookFunc
		afterResponseHookFuncs     []AfterResponseHookFunc
		beforeRequestArgsHookFuncs []BeforeRequestArgsHookFunc
		option                     []ModifySessionOption
		args                       SessionArgs
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
	url                string
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
		r.url = _url
	}
}

func Proxy(_proxy string) ModifySessionOption {
	return func(r *SessionArgs) {
		r.proxy = _proxy
	}
}

func Cookies(_cookies []*http.Cookie) ModifySessionOption {
	return func(r *SessionArgs) {
		r.cookies = _cookies
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
	args := SessionArgs{
		proxy:              "",
		timeout:            30 * time.Second,
		skipVerifyTLS:      false,
		disableKeepAlive:   false,
		disableCompression: false,
		allowRedirects:     true,
	}

	for _, f := range opts {
		f(&args)
	}

	tranSport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   args.timeout,
			KeepAlive: args.timeout,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: args.skipVerifyTLS,
		},
		DisableKeepAlives:  args.disableKeepAlive,
		DisableCompression: args.disableCompression,
	}
	if args.proxy != "" {
		Url, _ := url.Parse(args.proxy)
		proxyUrl := http.ProxyURL(Url)
		tranSport.Proxy = proxyUrl
	}

	client := &http.Client{}
	client.Transport = tranSport

	if args.allowRedirects {
		client.CheckRedirect = defaultCheckRedirect
	} else {
		client.CheckRedirect = disableRedirect
	}

	Url, _ := url.Parse(args.url)
	session := &Session{
		Url:       Url,
		option:    opts,
		CookieJar: NewCookieJar(),
		args:      args,
	}
	client.Jar = session.CookieJar
	session.Client = client

	if session.Url != nil && args.cookies != nil {
		session.CookieJar.SetCookies(session.Url, args.cookies)
	}
	return session
}

func (s *Session) SetUrl(_url string) *Session {
	s.Url, _ = url.Parse(_url)
	s.args.url = _url
	return s
}

func (s *Session) GetUrl() string {
	return s.args.url
}

func (s *Session) SetCookies(_url string, cookies []*http.Cookie) *Session {
	Url, err := url.Parse(_url)
	if err != nil {
		panic("set cookies error: " + err.Error())
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
	s.args.timeout = timeout
	return s
}

func (s *Session) GetTimeout() time.Duration {
	return s.args.timeout
}

func (s *Session) SetSkipVerifyTLS(skip bool) *Session {
	s.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = skip
	s.args.skipVerifyTLS = skip
	return s
}

func (s *Session) GetSkipVerifyTLS() bool {
	return s.args.skipVerifyTLS
}

func (s *Session) SetProxy(proxy string) *Session {
	proxyUrl, _ := url.Parse(proxy)
	s.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
	s.args.proxy = proxy
	return s
}

func (s *Session) GetProxy() string {
	return s.args.proxy
}

func (s *Session) SetDisableKeepAlive(disable bool) *Session {
	s.Client.Transport.(*http.Transport).DisableKeepAlives = disable
	return s
}

func (s *Session) GetDisableKeepAlive() bool {
	return s.args.disableKeepAlive
}

func (s *Session) SetDisableCompression(disable bool) *Session {
	s.Client.Transport.(*http.Transport).DisableCompression = disable
	return s
}

func (s *Session) GetDisableCompression() bool {
	return s.args.disableCompression
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
	return s.args.allowRedirects
}

func (s *Session) Save(path string, _url string) error {
	return s.CookieJar.Save(path, _url)
}

func (s *Session) Load(path string, _url string) error {
	return s.CookieJar.Load(path, _url)
}

func (s *Session) Request(method string, urlStr string, args RequestArgs) *Response {
	s.Lock()
	defer s.Unlock()

	method = strings.ToUpper(method)
	switch method {
	case HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH:

		for _, fn := range s.beforeRequestArgsHookFuncs {
			err := fn(&args)
			if err != nil {
				break
			}
		}

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
		s.request.Header.Set("User-Agent", RandomUserAgent(nil))

		if s.Client == nil {
			s.CookieJar = NewCookieJar()
			s.Client = &http.Client{}
			s.Client.Jar = s.CookieJar
			s.Client.Transport = &http.Transport{}
		}

		if &args != nil {
			err = args.setRequestOpt(s.request)
			if err != nil {
				return &Response{
					Err: err,
				}
			}

			err = args.setClientOpt(s.Client)
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

	resp := NewResponse(r)

	for _, fn := range s.afterResponseHookFuncs {
		err = fn(resp)
		if err != nil {
			break
		}
	}
	return resp
}

func (s *Session) AsyncRequest(method string, urlStr string, args RequestArgs, ch chan *Response) {
	response := s.Request(method, urlStr, args)
	ch <- response
}

func (s *Session) GetRequest() *http.Request {
	return s.request
}

type (
	BeforeRequestArgsHookFunc func(*RequestArgs) error
	BeforeRequestHookFunc     func(*http.Request) error
	AfterResponseHookFunc     func(*Response) error
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

func (s *Session) RegisterBeforeRequestArgsHook(fn BeforeRequestArgsHookFunc) error {
	if s.beforeRequestArgsHookFuncs == nil {
		s.beforeRequestArgsHookFuncs = make([]BeforeRequestArgsHookFunc, 0, 8)
	}
	if len(s.beforeRequestArgsHookFuncs) > 7 {
		return ErrHookFuncMaxLimit
	}
	s.beforeRequestArgsHookFuncs = append(s.beforeRequestArgsHookFuncs, fn)
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
	session.Client.Jar = s.CookieJar
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

func (s *Session) Get(url string, args RequestArgs) *Response {
	return s.Request("get", url, args)
}

func (s *Session) AsyncGet(url string, args RequestArgs, ch chan *Response) {
	go s.AsyncRequest("get", url, args, ch)
}

func (s *Session) Post(url string, args RequestArgs) *Response {
	return s.Request("post", url, args)
}

func (s *Session) AsyncPost(url string, args RequestArgs, ch chan *Response) {
	go s.AsyncRequest("post", url, args, ch)
}

func (s *Session) Head(url string, args RequestArgs) *Response {
	return s.Request("head", url, args)
}

func (s *Session) AsyncHead(url string, args RequestArgs, ch chan *Response) {
	go s.AsyncRequest("head", url, args, ch)
}

func (s *Session) Delete(url string, args RequestArgs) *Response {
	return s.Request("delete", url, args)
}

func (s *Session) AsyncDelete(url string, args RequestArgs, ch chan *Response) {
	go s.AsyncRequest("delete", url, args, ch)
}

func (s *Session) Options(url string, args RequestArgs) *Response {
	return s.Request("options", url, args)
}

func (s *Session) AsyncOptions(url string, args RequestArgs, ch chan *Response) {
	go s.AsyncRequest("options", url, args, ch)
}

func (s *Session) Put(url string, args RequestArgs) *Response {
	return s.Request("put", url, args)
}

func (s *Session) AsyncPut(url string, args RequestArgs, ch chan *Response) {
	go s.AsyncRequest("put", url, args, ch)
}

func (s *Session) Patch(url string, args RequestArgs) *Response {
	return s.Request("patch", url, args)
}

func (s *Session) AsyncPatch(url string, args RequestArgs, ch chan *Response) {
	go s.AsyncRequest("patch", url, args, ch)
}
