package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	ErrInvalidMethod = errors.New("requests: method is invalid")

	ErrFileInfo = errors.New("requests: invalid file information")

	ErrParamConflict = errors.New("requests: post conflict")

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

type Params map[string]string

func (d *Params) Update(s Params) {
	for k, v := range s {
		(*d)[k] = v
	}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

type Request struct {
	*http.Request
	Headers Headers
	Cookies SimpleCookie
	Auth    Auth
	Params  Params
	Form    Form
	Payload Payload
	Binary  []byte
	Files   []*File

	Proxy              string
	Timeout            time.Duration
	SkipVerifyTLS      bool
	Chunked            bool
	AllowRedirects     bool
	DisableKeepAlive   bool
	DisableCompression bool
}

func (r *Request) isConflict() bool {
	count := 0
	if r.Form != nil {
		count++
	}
	if r.Payload != nil {
		count++
	}
	if len(r.Files) != 0 {
		count++
	}
	if r.Binary != nil {
		count++
	}
	return count > 1
}

func (r *Request) setQuery() error {
	originURL := r.Request.URL
	extendQuery := make([]byte, 0)

	for k, v := range r.Params {
		kEscaped := url.QueryEscape(k)
		vEscaped := url.QueryEscape(v)

		extendQuery = append(extendQuery, '&')
		extendQuery = append(extendQuery, []byte(kEscaped)...)
		extendQuery = append(extendQuery, '=')
		extendQuery = append(extendQuery, []byte(vEscaped)...)
	}

	if originURL.RawQuery == "" {
		extendQuery = extendQuery[1:]
	}

	originURL.RawQuery += string(extendQuery)
	return nil
}

func (r *Request) setForm() error {
	data := ""
	for k, v := range r.Form {
		k = url.QueryEscape(k)

		vs, ok := v.(string)
		if !ok {
			return fmt.Errorf("post data %v[%T] must be string type", v, v)
		}
		vs = url.QueryEscape(vs)
		data = fmt.Sprintf("%s&%s=%s", data, k, vs)
	}

	data = data[1:]
	v := strings.NewReader(data)
	r.Request.Body = ioutil.NopCloser(v)
	r.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if !r.Chunked {
		r.Request.ContentLength = int64(v.Len())
	}
	return nil
}

func (r *Request) setPayload() error {
	jsonV, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}
	v := bytes.NewBuffer(jsonV)
	r.Request.Body = ioutil.NopCloser(v)
	r.Request.Header.Set("Content-Type", "application/json")
	if !r.Chunked {
		r.Request.ContentLength = int64(v.Len())
	}
	return nil
}

func (r *Request) setBinary() error {
	body := bytes.NewReader(r.Binary)
	r.Request.Body = ioutil.NopCloser(body)
	if !r.Chunked {
		r.Request.ContentLength = int64(len(r.Binary))
	}

	return nil
}

func (r *Request) setFiles() error {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	var fp *os.File
	defer func() {
		if fp != nil {
			fp.Close()
		}
	}()

	for _, file := range r.Files {
		mimeType := file.MimeType
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(file.Param), escapeQuotes(file.Name)))
		h.Set("Content-Type", mimeType)

		var fileWriter io.Writer
		var err error

		fileWriter, err = writer.CreatePart(h)
		if err != nil {
			return err
		}

		if len(file.Src) != 0 {
			_, err = fileWriter.Write(file.Src)
			if err != nil {
				return err
			}
		} else {
			fp, err = os.Open(file.Path)
			if err != nil {
				return err
			}

			_, err = io.Copy(fileWriter, fp)
			if err != nil {
				return err
			}
		}

		for name, value := range file.Args {
			err := writer.WriteField(name, value)
			if err != nil {
				return err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	r.Request.Body = ioutil.NopCloser(buffer)
	contentType := writer.FormDataContentType()
	r.Request.Header.Set("Content-Type", contentType)
	if !r.Chunked {
		r.Request.ContentLength = int64(buffer.Len())
	}
	return nil
}

func (r *Request) setAuth() error {
	for k, v := range r.Auth {
		vs, ok := v.(string)
		if !ok {
			return fmt.Errorf("basic-auth %v[%T] must be string type", v, v)
		}
		r.Request.SetBasicAuth(k, vs)
	}
	return nil
}

func (r *Request) setRequestArgs() error {
	if r.isConflict() {
		return ErrParamConflict
	}

	if r.Headers != nil {
		for key, value := range r.Headers {
			r.Request.Header[key] = []string{value}
		}
	}

	if r.Cookies != nil {
		for cookieK, cookieV := range r.Cookies {
			c := &http.Cookie{
				Name:  cookieK,
				Value: cookieV,
			}
			r.Request.AddCookie(c)
		}
	}

	if r.Auth != nil {
		err := r.setAuth()
		if err != nil {
			return err
		}
	}

	if r.Params != nil {
		err := r.setQuery()
		if err != nil {
			return err
		}
	}

	if r.Form != nil {
		err := r.setForm()
		if err != nil {
			return err
		}
	}

	if r.Payload != nil {
		err := r.setPayload()
		if err != nil {
			return err
		}
	}

	if r.Binary != nil {
		err := r.setBinary()
		if err != nil {
			return err
		}
	}

	if len(r.Files) != 0 {
		err := r.setFiles()
		if err != nil {
			return err
		}
	}

	return nil
}

func Get(url string, args ...interface{}) *Response {
	session := NewSession()
	return session.Get(url, args...)
}

func Post(url string, args ...interface{}) *Response {
	session := NewSession()
	return session.Post(url, args...)
}

func Head(url string, args ...interface{}) *Response {
	session := NewSession()
	return session.Head(url, args...)
}

func Delete(url string, args ...interface{}) *Response {
	session := NewSession()
	return session.Delete(url, args...)
}

func Options(url string, args ...interface{}) *Response {
	session := NewSession()
	return session.Options(url, args...)
}

func Put(url string, args ...interface{}) *Response {
	session := NewSession()
	return session.Put(url, args...)
}

func Patch(url string, args ...interface{}) *Response {
	session := NewSession()
	return session.Patch(url, args...)
}
