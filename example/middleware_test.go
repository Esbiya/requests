package example

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Esbiya/requests"
	"github.com/forgoer/openssl"
	"github.com/gofrs/uuid"
)

func genKey() string {
	UUID, _ := uuid.NewV4()
	return strings.Replace(UUID.String(), "-", "", -1)[:24]
}

func convertToBin(n int, bin int) string {
	var b string
	switch {
	case n == 0:
		for i := 0; i < bin; i++ {
			b += "0"
		}
	case n > 0:
		for ; n > 0; n /= 2 {
			b = strconv.Itoa(n%2) + b
		}
		j := bin - len(b)
		for i := 0; i < j; i++ {
			b = "0" + b
		}
	case n < 0:
		n = n * -1
		s := convertToBin(n, bin)
		for i := 0; i < len(s); i++ {
			if s[i:i+1] == "1" {
				b += "0"
			} else {
				b += "1"
			}
		}
		n, err := strconv.ParseInt(b, 2, 64)
		if err != nil {
			fmt.Println(err)
		}
		b = convertToBin(int(n+1), bin)
	}
	return b
}

func genCipher(data string) string {
	result := ""
	for i := 0; i < len(data); i++ {
		j := convertToBin(int(data[i]), 2)
		result += j + " "
	}
	return result[:len(result)-1]
}

func TestRequest(t *testing.T) {
	url := "http://wenshuapp.court.gov.cn/appinterface/rest.q4w"

	session := requests.NewSession()
	headers := requests.Headers{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "Dalvik/2.1.0 (Linux; U; Android 8.0.0; Google Nexus 5X Build/OPR6.170623.017)",
		"Host":         "wenshuapp.court.gov.cn",
	}
	now := time.Now()
	UUID, _ := uuid.NewV4()

	key := genKey()
	fmt.Println(key)
	iv := fmt.Sprintf("%s%s%s", now.Format("2006"), now.Format("01"), now.Format("02"))
	fmt.Println(iv)
	timestamp := strconv.FormatInt(now.UnixNano()/1e6, 10)

	_ = session.RegisterBeforeReqHook(func(req *requests.Request) error {
		encryptStr, err := openssl.Des3CBCEncrypt([]byte(timestamp), []byte(key), []byte(iv), openssl.PKCS7_PADDING)
		if err != nil {
			return err
		}
		(*req.Form)["params"].(map[string]interface{})["ciphertext"] = genCipher(key + iv + base64.StdEncoding.EncodeToString(encryptStr))

		d, err := json.Marshal(req.Form)
		fmt.Println(string(d))
		req.Form = &requests.Form{
			"request": base64.StdEncoding.EncodeToString(d),
		}
		return err
	})
	_ = session.RegisterAfterRespHook(func(response *requests.Response) error {
		result, err := response.JSON()
		if err != nil {
			return err
		}
		if result.Get("data.secretKey").Exists() {
			b, _ := base64.StdEncoding.DecodeString(result.Get("data.content").String())
			response.Bytes, err = openssl.Des3CBCDecrypt(b, []byte(result.Get("data.secretKey").String()), []byte(iv), openssl.PKCS7_PADDING)
			if err != nil {
				return err
			}
			response.Text = string(response.Bytes)
		}
		return err
	})

	data := requests.Form{
		"id":      fmt.Sprintf("%s%s%s%s%s%s", now.Format("2006"), now.Format("01"), now.Format("02"), now.Format("15"), now.Format("04"), now.Format("05")),
		"command": "queryDoc",
		"params": map[string]interface{}{
			"devid":    strings.Replace(UUID.String(), "-", "", -1),
			"devtype":  "1",
			"pageSize": "20", "sortFields": "s50:desc", "pageNum": "1",
			"queryCondition": []map[string]interface{}{
				{"key": "s2", "value": "四川省成都市中级人民法院"},
			},
		},
	}
	cookies := requests.SimpleCookie{
		"SESSION": "b52c07c2-d1e7-4a65-9a80-d80e181bf4b8",
	}
	resp := session.Post(url, headers, data, cookies, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true})
	if resp.Error() != nil {
		panic(resp.Error())
	}
	result, _ := resp.JSON()
	xx, _ := json.MarshalIndent(result, "", "    ")
	log.Println(string(xx))
	log.Println(resp.Cost().String())
}
