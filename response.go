package requests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/axgle/mahonia"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	*http.Response
	encoding string
	cost     time.Duration
	Text     string
	Bytes    []byte
	err      error
}

func NewResponse(r *http.Response, cost time.Duration) *Response {
	resp := &Response{
		Response: r,
		encoding: "utf-8",
		cost:     cost,
		Text:     "",
		Bytes:    []byte{},
	}

	err := resp.bytes()
	resp.err = err
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

func (r *Response) Error() error {
	return r.err
}

func (r *Response) Cost() time.Duration {
	return r.cost
}

func (r *Response) XML(v interface{}) error {
	return xml.Unmarshal(r.Bytes, v)
}

func (r *Response) JSON() (map[string]interface{}, error) {
	if r.err != nil {
		return nil, r.err
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
	if r.err != nil {
		return nil, r.err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response code: " + strconv.Itoa(r.StatusCode))
	}
	var result map[string]interface{}
	re, _ := regexp.Compile("\\({[\\s\\S]*?}\\)")
	y := re.FindStringSubmatch(r.Text)
	if len(y) == 0 {
		return result, ErrNotJSONResponse
	}
	decoder := json.NewDecoder(bytes.NewReader([]byte(y[0][1 : len(y[0])-1])))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		return result, ErrNotJSONResponse
	}
	return result, nil
}

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

func (r Response) GetEncode() string {
	return r.encoding
}

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
