package requests

import (
	"bytes"
	"encoding/json"
	"github.com/axgle/mahonia"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Response is the wrapper for http.Response
type Response struct {
	*http.Response
	encoding string
	Status   int
	Text     string
	Bytes    []byte
	Err      error
}

func NewResponse(r *http.Response) *Response {
	resp := &Response{
		Response: r,
		encoding: "utf-8",
		Text:     "",
		Status:   r.StatusCode,
		Bytes:    []byte{},
	}

	err := resp.bytes()
	resp.Err = err
	resp.text()
	return resp
}

func (r *Response) text() {
	r.Text = string(r.Bytes)
}

func (r *Response) bytes() error {
	defer r.Body.Close()

	var err error
	r.Bytes, err = ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(r.Bytes))
	return err
}

func (r *Response) JSON() (map[string]interface{}, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response code: " + strconv.Itoa(r.StatusCode))
	}
	var result map[string]interface{}
	dec := json.NewDecoder(bytes.NewBuffer(r.Bytes))
	// 将处理的数字转化成 json.Number 的形式，防止丢失精度
	dec.UseNumber()
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Response) CallbackJSON() (map[string]interface{}, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response code: " + strconv.Itoa(r.StatusCode))
	}
	var result map[string]interface{}
	re, _ := regexp.Compile("\\({[\\s\\S]*?}\\)")
	y := re.FindStringSubmatch(r.Text)
	decoder := json.NewDecoder(bytes.NewReader([]byte(y[0][1 : len(y[0])-1])))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// SetEncode changes Response.encoding
// and it changes Response.Text every times be invoked
func (r *Response) SetEncode(e string) error {
	if r.encoding != e {
		r.encoding = strings.ToLower(e)
		decoder := mahonia.NewDecoder(e)
		if decoder == nil {
			return ErrUnrecognizedEncoding
		}
		r.Text = decoder.ConvertString(r.Text)
	}
	return nil
}

// GetEncode returns Response.encoding
func (r Response) GetEncode() string {
	return r.encoding
}

// SaveFile save bytes data to a local file
func (r Response) SaveFile(filename string) error {
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.Write(r.Bytes)
	if err != nil {
		return err
	}
	return nil
}
