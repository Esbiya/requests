package requests

import (
	"encoding/json"
	"github.com/wuxiaoxiaoshen/fakeuseragent/application"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	Chrome = 1
	Safari = 2
	IE     = 3
	Opera  = 4
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func RandomNum(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

func RandomUserAgent(t interface{}) string {
	fakeUserAgent := application.NewFakeUserAgent(true, false, true)
	if t == nil {
		t = RandomNum(1, 4)
	}
	switch t {
	case Chrome:
		return fakeUserAgent.Chrome()
	case Safari:
		return fakeUserAgent.Safari()
	case IE:
		return fakeUserAgent.IE()
	case Opera:
		return fakeUserAgent.Opera()
	default:
		return fakeUserAgent.Random()
	}
}

func TransferCookie(c map[string]interface{}) (*http.Cookie, error) {
	var cookie *http.Cookie
	cBytes, err := json.Marshal(c)
	if err != nil {
		return cookie, err
	}
	err = json.Unmarshal(cBytes, &cookie)
	return cookie, err
}

func TransferCookies(c []map[string]interface{}) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	var err error
	for _, c1 := range c {
		var c2 *http.Cookie
		c2, err = TransferCookie(c1)
		cookies = append(cookies, c2)
	}
	return cookies, err
}
