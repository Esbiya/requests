package requests

import (
	"github.com/wuxiaoxiaoshen/fakeuseragent/application"
	"math/rand"
	"os"
	"strconv"
	"time"
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
	case 1:
		return fakeUserAgent.Chrome()
	case 2:
		return fakeUserAgent.Safari()
	case 3:
		return fakeUserAgent.IE()
	case 4:
		return fakeUserAgent.Opera()
	default:
		return fakeUserAgent.Random()
	}
}
