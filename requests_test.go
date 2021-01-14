package requests

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"reflect"
	"sort"
	"testing"
)

func EachMap(eachMap interface{}, eachFunc interface{}) {
	eachMapValue := reflect.ValueOf(eachMap)
	eachFuncValue := reflect.ValueOf(eachFunc)
	eachMapType := eachMapValue.Type()
	eachFuncType := eachFuncValue.Type()
	if eachMapValue.Kind() != reflect.Map {
		panic(errors.New("ksort.EachMap failed. parameter \"eachMap\" type must is map[...]...{}"))
	}
	if eachFuncValue.Kind() != reflect.Func {
		panic(errors.New("ksort.EachMap failed. parameter \"eachFunc\" type must is func(key ..., value ...)"))
	}
	if eachFuncType.NumIn() != 2 {
		panic(errors.New("ksort.EachMap failed. \"eachFunc\" input parameter count must is 2"))
	}
	if eachFuncType.In(0).Kind() != eachMapType.Key().Kind() {
		panic(errors.New("ksort.EachMap failed. \"eachFunc\" input parameter 1 type not equal of \"eachMap\" key"))
	}
	if eachFuncType.In(1).Kind() != eachMapType.Elem().Kind() {
		panic(errors.New("ksort.EachMap failed. \"eachFunc\" input parameter 2 type not equal of \"eachMap\" value"))
	}

	// 对 key 进行排序
	// 获取排序后 map 的 key 和 value，作为参数调用 eachFunc 即可
	switch eachMapType.Key().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		keys := make([]int, 0)
		keysMap := map[int]reflect.Value{}
		for _, value := range eachMapValue.MapKeys() {
			keys = append(keys, int(value.Int()))
			keysMap[int(value.Int())] = value
		}
		sort.Ints(keys)
		for _, key := range keys {
			eachFuncValue.Call([]reflect.Value{keysMap[key], eachMapValue.MapIndex(keysMap[key])})
		}
	case reflect.Float64, reflect.Float32:
		keys := make([]float64, 0)
		keysMap := map[float64]reflect.Value{}
		for _, value := range eachMapValue.MapKeys() {
			keys = append(keys, value.Float())
			keysMap[value.Float()] = value
		}
		sort.Float64s(keys)
		for _, key := range keys {
			eachFuncValue.Call([]reflect.Value{keysMap[key], eachMapValue.MapIndex(keysMap[key])})
		}
	case reflect.String:
		keys := make([]string, 0)
		keysMap := map[string]reflect.Value{}
		for _, value := range eachMapValue.MapKeys() {
			keys = append(keys, value.String())
			keysMap[value.String()] = value
		}
		sort.Strings(keys)
		for _, key := range keys {
			eachFuncValue.Call([]reflect.Value{keysMap[key], eachMapValue.MapIndex(keysMap[key])})
		}
	default:
		panic(errors.New("\"eachMap\" key type must is int or float or string"))
	}
}

func TestRequest(t *testing.T) {
	url := "https://www.baidu.com/"

	session := NewSession()

	_ = session.RegisterBeforeRequestArgsHook(func(args *RequestArgs) error {
		args.Proxy = "http://127.0.0.1:8888"
		args.SkipVerifyTLS = false
		// 对 params 进行排序拼接 base64 加签
		signature := ""
		EachMap(args.Params, func(key string, value interface{}) {
			signature += key + "=" + value.(string)
		})
		args.Params["signature"] = base64.StdEncoding.EncodeToString([]byte(signature))
		return nil
	})
	_ = session.RegisterBeforeReqHook(func(request *http.Request) error {
		request.Close = true
		return nil
	})
	_ = session.RegisterAfterRespHook(func(response *Response) error {
		var err error
		response.Bytes, err = base64.StdEncoding.DecodeString(response.Text)
		response.Text = string(response.Bytes)
		return err
	})

	params := DataMap{
		"xxx": "222",
		"aaa": "heheh",
	}
	resp := session.Get(url, RequestArgs{SkipVerifyTLS: true, Params: params})
	log.Println(resp.Text)

	c := make(chan *Response, 1)
	session.AsyncGet(url, RequestArgs{}, c)
	log.Println((<-c).JSON())
}
