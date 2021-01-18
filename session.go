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
		args                   SessionArgs
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

func (h *SessionArgs) setClientOpt(client *http.Client) error {
	if !h.AllowRedirects {
		client.CheckRedirect = disableRedirect
	}

	client.Timeout = h.Timeout * time.Second

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

type ModifySessionOption func(session *SessionArgs)

func SetUrl(_url string) ModifySessionOption {
	return func(r *SessionArgs) {
		r.Url = _url
	}
}

func SetProxy(_proxy string) ModifySessionOption {
	return func(r *SessionArgs) {
		r.Proxy = _proxy
	}
}

func SetCookies(_cookies []*http.Cookie) ModifySessionOption {
	return func(r *SessionArgs) {
		r.Cookies = _cookies
	}
}

func SetTimeout(_timeout time.Duration) ModifySessionOption {
	return func(r *SessionArgs) {
		r.Timeout = _timeout
	}
}

func SetSkipVerifyTLS(_skipVerifyTLS bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.SkipVerifyTLS = _skipVerifyTLS
	}
}

func SetDisableKeepAlive(_disableKeepAlive bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.DisableKeepAlive = _disableKeepAlive
	}
}

func SetDisableCompression(_disableCompression bool) ModifySessionOption {
	return func(r *SessionArgs) {
		r.DisableCompression = _disableCompression
	}
}

func NewSession(opts ...ModifySessionOption) *Session {
	v := SessionArgs{
		Proxy:              "",
		Timeout:            30 * time.Second,
		SkipVerifyTLS:      false,
		DisableKeepAlive:   false,
		DisableCompression: false,
		AllowRedirects:     true,
	}

	for _, f := range opts {
		f(&v)
	}

	tranSport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   v.Timeout,
			KeepAlive: v.Timeout,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: v.SkipVerifyTLS,
		},
		DisableKeepAlives:     v.DisableKeepAlive,
		DisableCompression:    v.DisableCompression,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if v.Proxy != "" {
		Url, _ := url.Parse(v.Proxy)
		proxyUrl := http.ProxyURL(Url)
		tranSport.Proxy = proxyUrl
	}

	client := &http.Client{}
	client.Transport = tranSport

	if v.AllowRedirects {
		client.CheckRedirect = defaultCheckRedirect
	} else {
		client.CheckRedirect = disableRedirect
	}

	Url, _ := url.Parse(v.Url)
	session := &Session{
		Url:       Url,
		option:    opts,
		CookieJar: NewCookieJar(),
		request:   &Request{},
		args:      v,
	}
	client.Jar = session.CookieJar
	session.Client = client

	if session.Url != nil && v.Cookies != nil {
		session.CookieJar.SetCookies(session.Url, v.Cookies)
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
	resp := s.Do(method, urlStr, args)
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
				Err: err,
			}
		}
		urlStrParsed.RawQuery = urlStrParsed.Query().Encode()

		req, err := http.NewRequest(method, urlStrParsed.String(), nil)
		if err != nil {
			return &Response{
				Err: err,
			}
		}
		s.request.Request = req
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
			case Files:
				s.request.Files = _arg
			case http.Header:
				for key, values := range _arg {
					for _, value := range values {
						s.request.Header.Add(key, value)
					}
				}
			case string:
				body := strings.NewReader(_arg)
				s.request.Body = ioutil.NopCloser(body)
				s.request.ContentLength = int64(len(_arg))
			case []byte:
				body := bytes.NewReader(_arg)
				s.request.Body = ioutil.NopCloser(body)
				s.request.ContentLength = int64(len(_arg))
			case *http.Cookie:
				s.request.AddCookie(_arg)
			case []*http.Cookie:
				for _, cookie := range _arg {
					s.request.AddCookie(cookie)
				}
			case SessionArgs:
				err := _arg.setClientOpt(s.Client)
				if err != nil {
					return &Response{
						Err: err,
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

		err = s.request.setRequestOpt()
		if err != nil {
			return &Response{
				Err: err,
			}
		}

	default:
		return &Response{
			Err: ErrInvalidMethod,
		}
	}

	r, err := s.Client.Do(s.request.Request)
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

func (s *Session) AsyncGet(url string, c chan *Response, args ...interface{}) {
	go s.AsyncDo("get", url, c, args)
}

func (s *Session) Post(url string, args ...interface{}) *Response {
	return s.Do("post", url, args...)
}

func (s *Session) AsyncPost(url string, c chan *Response, args ...interface{}) {
	go s.AsyncDo("post", url, c, args)
}

func (s *Session) Head(url string, args ...interface{}) *Response {
	return s.Do("head", url, args...)
}

func (s *Session) AsyncHead(url string, c chan *Response, args ...interface{}) {
	go s.AsyncDo("head", url, c, args)
}

func (s *Session) Delete(url string, args ...interface{}) *Response {
	return s.Do("delete", url, args...)
}

func (s *Session) AsyncDelete(url string, c chan *Response, args ...interface{}) {
	go s.AsyncDo("delete", url, c, args)
}

func (s *Session) Options(url string, args ...interface{}) *Response {
	return s.Do("options", url, args...)
}

func (s *Session) AsyncOption(url string, c chan *Response, args ...interface{}) {
	go s.AsyncDo("option", url, c, args)
}

func (s *Session) Put(url string, args ...interface{}) *Response {
	return s.Do("put", url, args...)
}

func (s *Session) AsyncPut(url string, c chan *Response, args ...interface{}) {
	go s.AsyncDo("put", url, c, args)
}

func (s *Session) Patch(url string, args ...interface{}) *Response {
	return s.Do("patch", url, args...)
}

func (s *Session) AsyncPatch(url string, c chan *Response, args ...interface{}) {
	go s.AsyncDo("patch", url, c, args)
}
