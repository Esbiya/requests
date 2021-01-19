## requests

​		一个东拼西凑乱七八糟性能极差的符合自己使用习惯的类似 **`python requests`** 的 **`golang http`** 请求库，**`请勿使用`**。

## 安装

```shell
go get -u github.com/Esbiya/requests
```

## 使用

```go
package main

import (
    "github.com/Esbiya/requests"
    "log"
    "net/http"
)

func main() {
	url := "https://www.baidu.com/"
  	resp := requests.Get(url)
	  if resp.StatusCode != http.StatusOK {
		    log.Fatal("状态码异常")
   	}
   	log.Println(resp.Text)
}
```

### get 请求

```go
params := requests.Params{"1": "2"}
resp := requests.Get("https://www.baidu.com", params, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true})
if resp.Error() != nil {
	log.Fatal(resp.Error())
}
log.Println(resp.Text)
log.Println(resp.Cost().String())
```

### 异步 get 请求

```go
start := time.Now()
for i := 0; i < 10; i++ {
	headers := requests.Headers{
		"Connection": "keep-alive",
	}
	requests.AsyncGet("https://www.baidu.com", headers).Then(func(r *requests.Response) {
	    if r.Error() != nil {
	    	log.Fatal(r.Error())
		}
		log.Println("function cost => " + r.Cost().String())
	})
}
log.Println("execute first")
requests.AsyncWait() // 必须等待, 否则异步请求不会执行
end := time.Now()
log.Println("all cost => " + end.Sub(start).String())
```

### post form 请求

```go
url := "https://accounts.douban.com/j/mobile/login/basic"
headers := requests.Headers{
	"Connection":       "keep-alive",
	"Pragma":           "no-cache",
	"Cache-Control":    "no-cache",
	"Accept":           "application/json",
	"X-Requested-With": "XMLHttpRequest",
	"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
	"Content-Type":     "application/x-www-form-urlencoded",
	"Origin":           "https://accounts.douban.com",
	"Referer":          "https://accounts.douban.com/passport/login_popup?login_source=anony",
	"Accept-Language":  "zh-CN,zh;q=0.9",
}
cookies := requests.SimpleCookie{
	"bid":               "lj_mS940akg",
	"douban-fav-remind": "1",
	"ll":                "118281",
	"_vwo_uuid_v2":      "D80341DC04F297D12F96A751576B82F67|271ec79a2a95920f60676d36ea786643",
	"__gads":            "ID=0930c55e6f0d2ebe-22d4125465c400e7:T=1603714878:RT=1603714878:S=ALNI_Ma45i6lqTVh5ATe4v0iNBhhBVkUJQ",
	"__utmc":            "30149280",
	"apiKey":            "",
	"__utma":            "30149280.1034523831.1589204237.1609676365.1611036335.28",
	"__utmz":            "30149280.1611036335.28.27.utmcsr=baidu|utmccn=(organic)|utmcmd=organic",
	"__utmt":            "1",
	"__utmb":            "30149280.1.10.1611036335",
	"login_start_time":  "1611036340462",
}
data := requests.Form{
	"ck":       "",
	"remember": "true",
	"name":     "18829040039",
	"password": "wefewfw",
}
resp := requests.Post(url, headers, data, cookies, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true})
if resp.Error() != nil {
	log.Fatal(resp.Error())
}
log.Println(resp.Text)
```

### 异步 post 请求

```go
url := "https://accounts.douban.com/j/mobile/login/basic"
headers := requests.Headers{
	"Connection":       "keep-alive",
	"Pragma":           "no-cache",
	"Cache-Control":    "no-cache",
	"Accept":           "application/json",
	"X-Requested-With": "XMLHttpRequest",
	"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
	"Content-Type":     "application/x-www-form-urlencoded",
	"Origin":           "https://accounts.douban.com",
	"Referer":          "https://accounts.douban.com/passport/login_popup?login_source=anony",
	"Accept-Language":  "zh-CN,zh;q=0.9",
}
cookies := requests.SimpleCookie{
	"bid":               "lj_mS940akg",
	"douban-fav-remind": "1",
	"ll":                "118281",
	"_vwo_uuid_v2":      "D80341DC04F297D12F96A751576B82F67|271ec79a2a95920f60676d36ea786643",
	"__gads":            "ID=0930c55e6f0d2ebe-22d4125465c400e7:T=1603714878:RT=1603714878:S=ALNI_Ma45i6lqTVh5ATe4v0iNBhhBVkUJQ",
	"__utmc":            "30149280",
	"apiKey":            "",
	"__utma":            "30149280.1034523831.1589204237.1609676365.1611036335.28",
	"__utmz":            "30149280.1611036335.28.27.utmcsr=baidu|utmccn=(organic)|utmcmd=organic",
	"__utmt":            "1",
	"__utmb":            "30149280.1.10.1611036335",
	"login_start_time":  "1611036340462",
}
data := requests.Form{
	"ck":       "",
	"remember": "true",
	"name":     "18829040039",
	"password": "wefewfw",
}
requests.AsyncPost(url, headers, data, cookies, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true}).Then(func(r *requests.Response) {
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	log.Println(r.Text)
})
requests.AsyncWait()
```

### post payload 请求

```go
url := "https://www.guilinbank.com.cn/api-portal/portal-home/article/manage/queryListByLibraryId"
headers := requests.Headers{
	"Connection":       "keep-alive",
	"Accept":           "*/*",
	"X-Requested-With": "XMLHttpRequest",
	"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
	"Content-Type":     "application/json; charset=UTF-8",
	"Origin":           "https://www.guilinbank.com.cn",
	"Referer":          "https://www.guilinbank.com.cn/page-adapt/index/docList?libraryId=LIBRARY_productNotice&title=%E4%BA%A7%E5%93%81%E5%85%AC%E5%91%8A",
	"Accept-Language":  "zh-CN,zh;q=0.9",
}

data := requests.Payload{
	"pageSize":  15,
	"pageNum":   1,
	"title":     "",
	"endTime":   nil,
	"startTime": nil,
	"libraryId": "LIBRARY_productNotice$guilin_bank_portal",
}
// 等同于 json 字符串, == json.dumps({"1": "2"}) (python)
resp := requests.Post(url, headers, data, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true, DisableCompression: false})
if resp.Error() != nil {
	log.Fatal(resp.Error())
}
log.Println(resp.Text)
```

### post binary

```go
api := "http://192.168.100.107:7788"
data, err := ioutil.ReadFile("captcha.png")
if err != nil {		
    log.Fatal(err)
}
resp := requests.Post(api, data)
if resp.Error() != nil {
	log.Fatal(resp.Error())
}
log.Println(resp.Text)
```

### cookie 设置

```go
url := "https://www.baidu.com/"
data := requests.Data{
	"hello": "there",
}
headers := requests.Headers{
	"User-Agent": "xxx",
	"Cookie": "111=222; 333=444",  // 设置方式 1
}
Cookies := requests.SimpleCookies{
	"111": "222", 
    "333": "444",
}
resp := requests.Post(url, headers, data, cookie)
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 代理设置

- 支持 s5/http/https

```go
url := "https://www.baidu.com/"
resp := requests.Post(url, requests.Arguments{Proxy: "http://127.0.0.1:8888"})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 禁止重定向

```go
url := "https://www.baidu.com/"
resp := requests.Post(url, requests.Arguments{AllowRedirects: false})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 禁止响应数据压缩

```go
url := "https://www.baidu.com/"
resp := requests.Post(url, requests.Arguments{DisableCompression: true})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 禁止 tcp 链接复用

```go
url := "https://www.baidu.com/"
resp := requests.Post(url, requests.Arguments{DisableKeepAlive: true})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 上传文件

```go
session := requests.NewSession().SetProxy("http://127.0.0.1:8888").SetSkipVerifyTLS(true)
_ = session.Load("test.json", "https://chat.zalo.me/")
headers := requests.Headers{	
	"Accept":          "application/json, text/plain, */*",
	"Referer":         "https://chat.zalo.me/",
	"Accept-Language": "zh-CN,zh;q=0.9",
}
params := requests.Params{
	"zpw_ver":  "72",
	"zpw_type": "30",
	"params":   "CqJxEQSsxVJy3Rulu7wzpun/vVJeuzbZacam2v9JbiI=",
}
// 方式 1
file := requests.FileFromPath("/Users/esbiya/Desktop/pythonProjects/gitlab/whatsapp/test-0.png")
file.Param = "fileContent"
// 方式 2
b, _ := ioutil.ReadFile("/Users/esbiya/Desktop/pythonProjects/gitlab/whatsapp/test-0.png")
file1 := requests.FileFromBytes("test.png", b)
file1.Param = "fileContent"
// 方式 3, 细粒度配置
file2 := &requests.File{
	Name: "test.png",     													  // 文件名
	Param: "fileContent",  												      // 文件参数名
	Path: "/Users/esbiya/Desktop/pythonProjects/gitlab/whatsapp/test-0.png",  // 文件路径
	Src: []byte{},               										      // 文件内容
	Args: map[string]string{},   											  // 其他参数
}
// 支持多表单上传
resp := session.Post("https://files-wpa.chat.zalo.me/api/message/upthumb", headers, params, file, file1, file2)
if resp.Error() != nil {
	log.Fatal(resp.Error())
}
log.Println(resp.Text)
log.Println(resp.Cost().String())
```

### session

```go
session := requests.NewSession()
```

### session 设置

```go
session := requests.NewSession(
	requests.Url("http://www.baidu.com")  // 设置 cookies 全局 url
    requests.Cookies([]map[string]interface{}{"1": "2"}),  // 设置 cookies
    requests.Proxy("http://127.0.0.1:8888"),  // 设置代理
    requests.Timeout(time.Duration(5) * time.Second),  // 设置请求超时
    requests.SkipVerifyTLS: true,   // 忽略证书验证
    requests.Chunked(true),     // 设置是否分段上传
    requests.AllowRedirects(false),    // 是否禁止重定向
    requests.DisableKeepAlive(true),    // 禁止 tcp 连接复用
    requests.DisableCompression(true),  // 禁止响应数据压缩
)
```

或者:

```go
session := requests.NewSession().
	SetUrl("http://www.baidu.com").
	SetCookies([]map[string]interface{}{"1": "2"}).
	SetProxy("http://127.0.0.1:8888").
	SetTimeout(time.Duration(5) * time.Second).
	SetSkipVerifyTLS(true).
	SetChunked(true).
	SetAllowRedirects(false).
	SetDisableKeepAlive(true).
	SetDisableCompression(true)
```

### session get 请求

```go
resp := session.Get("http://www.baidu.com/")
```

### 异步 session get 请求

```go
start := time.Now()
session := requests.NewSession()
for i := 0; i < 10; i++ {
	c := make(chan *requests.Response, 1)
	headers := requests.Headers{
		"Connection": "keep-alive",
	}
	session.AsyncGet("https://www.baidu.com", c, headers).Then(func(r *requests.Response) {
		if r.Error() != nil {
			log.Fatal(r.Error())
		}
		log.Println("function cost => " + r.Cost().String())
	})
}
log.Println("execute first")
requests.AsyncWait() // 必须等待, 否则异步请求不会执行
end := time.Now()
log.Println("all cost => " + end.Sub(start).String())
```

### session post 请求

```go
resp := session.Post("http://www.baidu.com/")
```

### 异步 session post 请求

```go
url := "https://accounts.douban.com/j/mobile/login/basic"
headers := requests.Headers{
	"Connection":       "keep-alive",
	"Pragma":           "no-cache",
	"Cache-Control":    "no-cache",
	"Accept":           "application/json",
	"X-Requested-With": "XMLHttpRequest",
	"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
	"Content-Type":     "application/x-www-form-urlencoded",
	"Origin":           "https://accounts.douban.com",
	"Referer":          "https://accounts.douban.com/passport/login_popup?login_source=anony",
	"Accept-Language":  "zh-CN,zh;q=0.9",
}
cookies := requests.SimpleCookie{
	"bid":               "lj_mS940akg",
	"douban-fav-remind": "1",
	"ll":                "118281",
	"_vwo_uuid_v2":      "D80341DC04F297D12F96A751576B82F67|271ec79a2a95920f60676d36ea786643",
	"__gads":            "ID=0930c55e6f0d2ebe-22d4125465c400e7:T=1603714878:RT=1603714878:S=ALNI_Ma45i6lqTVh5ATe4v0iNBhhBVkUJQ",
	"__utmc":            "30149280",
	"apiKey":            "",
	"__utma":            "30149280.1034523831.1589204237.1609676365.1611036335.28",
	"__utmz":            "30149280.1611036335.28.27.utmcsr=baidu|utmccn=(organic)|utmcmd=organic",
	"__utmt":            "1",
	"__utmb":            "30149280.1.10.1611036335",
	"login_start_time":  "1611036340462",
}
data := requests.Form{
	"ck":       "",
	"remember": "true",
	"name":     "18829040039",
	"password": "wefewfw",
}
requests.AsyncPost(url, headers, data, cookies, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true}).Then(func(r *requests.Response) {
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	log.Println(r.Text)
})
requests.AsyncWait()
```

### 获取 session cookies

```go
// 获取指定域名下的 cookie
session.Cookies("http://jd.com/")
session.CookieJar.Get("https://jd.com/").String()   // 转化为字符串类型, 如 1=2; 3=4
session.CookieJar.Get("https://jd.com/").Map()      // 转化为 map 结构体, 如 map[string]interface{}{"1": "2"}
session.CookieJar.Get("https://jd.com/").Array()    // 转化为 []map[string]interface{}{"Name": "1", "Value": "2", "Domain": "jd.com"}

// 还可以用如下方式获取
session.CookieJar.String("https://jd.com/")
session.CookieJar.Map("https://jd.com/")
session.CookieJar.Array("https://jd.com/")
```

### 获取 session 其他设置

```go
url := session.GetUrl()
proxy := session.GetProxy()
timeout := session.GetTimeout()
ssl := session.GetSkipVerifyTLS()
allowRedirects := session.GetAllowRedirects()
disableKeepAlive := session.GetDisableKeepAlive()
disableCompression := session.GetDisableCompression()
```

### session 复制

```go
session1 := session.Copy()
```

### session 中间件注册

* 请求前对请求参数预处理, 如对请求参数进行排序加签操作

```go
_ = session.RegisterBeforeReqHook(func(req *requests.Request) error {
	encryptStr, err := openssl.Des3CBCEncrypt([]byte(timestamp), []byte(key), []byte(iv), openssl.PKCS7_PADDING)
	if err != nil {
		return err
	}
	req.Form["params"].(map[string]interface{})["ciphertext"] = genCipher(key + iv + base64.StdEncoding.EncodeToString(encryptStr))
	d, err := json.Marshal(req.Form)
	req.Form = requests.Form{
		"request": base64.StdEncoding.EncodeToString(d),
	}
	return err
})
```

* 请求完成对响应进行预处理, 如将加密响应解密成明文

```go
_ = session.RegisterAfterRespHook(func(response *requests.Response) error {
	result, err := response.JSON()
	if err != nil {
		return err
	}
	data := result["data"].(map[string]interface{})
	if _, ok := data["secretKey"]; ok {
		b, _ := base64.StdEncoding.DecodeString(data["content"].(string))
		response.Bytes, err = openssl.Des3CBCDecrypt(b, []byte(data["secretKey"].(string)), []byte(iv), openssl.PKCS7_PADDING)
		if err != nil {
			return err
		}
		response.Text = string(response.Bytes)
	}
	return err
})
```

### 响应

```go
resp := session.Post("http://www.baidu.com/")
resp.StatusCode          // 状态码
resp.Bytes               // 响应字节
resp.Text                // 响应字符串
resp.JSON()              // 解析 json 响应
resp.CallbackJSON()      // 解析 jQuery({"1": "2"}) json 响应
resp.SetEncode("gbk")    // 设置响应编码
resp.GetEncode()         // 获取响应编码
resp.SaveFile("xxx.jpg") // 文件写入
resp.Error()             // 响应错误, 正常响应为 nil
resp.Header              // 响应头
resp.Cookies()           // 响应 cookies, []*http.Cookie
resp.Location()          // 跳转 url
resp.ContentLength       // 响应内容大小
resp.Close               // tcp 连接是否已关闭 bool
reso.Cost()              // 请求耗时
```

### 代码自动生成

* 复制请求的 curl 使用脚本 generateScript.js 即可生成标准 requests 代码

```
console.log(curl2GoRequests(`curl 'https://login.taobao.com/newlogin/login.do?appName=taobao&fromSite=0&_bx-v=1.1.20' \\
  -H 'authority: login.taobao.com' \\
  -H 'pragma: no-cache' \\
  -H 'cache-control: no-cache' \\
  -H 'eagleeye-sessionid: t6kgpjLFz2m58n10zk21xtIrgabF' \\
  -H 'accept: application/json, text/plain, */*' \\
  -H 'eagleeye-pappname: gf3el0xc6g@256d85bbd150cf1' \\
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36' \\
  -H 'eagleeye-traceid: b216fbdc1610766945306100250cf1' \\
  -H 'content-type: application/x-www-form-urlencoded' \\
  -H 'origin: https://login.taobao.com' \\
  -H 'sec-fetch-site: same-origin' \\
  -H 'sec-fetch-mode: cors' \\
  -H 'sec-fetch-dest: empty' \\
  -H 'referer: https://login.taobao.com/member/login.jhtml?style=miniall&newMini2=true&full_redirect=true&&redirectURL=http%3A%2F%2Fworld.taobao.com%2F&from=worldlogin&minipara=1,1,1&from=worldlogin' \\
  -H 'accept-language: zh-CN,zh;q=0.9' \\
  -H 'cookie: tracknick=; miid=1928852900179986723; sgcookie=E1UIaaaplkSrC4cTJf4kP; _samesite_flag_=true; xlly_s=1; _m_h5_tk=7c054d74754411dc668bce986dd75216_1610775569864; _m_h5_tk_enc=55a0fdf4a40433f9ae9d0594c3c4f42e; cookie2=1c386caaf67de4d7690e183d5f775c4d; t=4e7cb12c7ff6bcb67d8b026c2a526d12; _tb_token_=e1d14f57e6ee3; cna=xC68F5h8QW8CATsq7e47o7pt; _bl_uid=88kaOj6wzI55Us131m251mUu1tCC; mt=ci=-1_0; _fbp=fb.1.1610766936473.1947231865; uc1=cookie14=Uoe1gq9r2Iujww%3D%3D; hng=CN%7Czh-CN%7CCNY%7C156; thw=cn; XSRF-TOKEN=a317e807-7ec1-4262-951b-5a24d963b790; tfstk=c7QdBombV5hKHLOOgMEiFyRrF-NcZdFpQXOKySbXJXHbW6uRiR7ckkg0ApkpW5C..; l=eBOYAwmHQG1jr4gvBOfZnurza779hIRcguPzaNbMiOCPO7fH5xMAW6G4BHLMCnGVnsEBr3zoKBo0B-LKZyUgl6Yl3ZQ7XPQoPdTh.; isg=BGpqw48L_idKrkx8cro3jJM0u9YM2-416jfPyfQjdL1IJwnh1muRRVGRt1M712bN' \\
  --data-raw 'loginId=18829040039&password2=84297912adeed3f348b3da7c2df2cec540f37423d748f3323f8e1e183848e41d63fb5f5bb222a3e333370c0c454fc364b06ab588c9d9a4ece3879d95f30d62ea9f87384edf16242502e457cbe52b12e7d7afce0b45033ff2cc1864ed40afccdcfd3eff46efb3d933fdc5f1d00872458ee9586c2dfb4fbfd7ec36f2cceab3fb6b&keepLogin=false&ua=140%23PrXDuQbczzFxwQo22Z3TCtSdvWMkKWzgkCqlMM6vlUV6B1j5G5TbbT8jtnpqJRQFCOSTi2n4HZSaCE3n5gycrvWqlbzx6IUpctgqzzr2hX0iU6OzzPzbVXlqlbrrNwQ1txyWEAziL2ILllfzzPzibvjj0TTd2eDKQt8NzeOVu2Ygl%2BFozD3RrfV9ONdOHaU3%2F0ttkrcT75HE391Y8Wc6THI4ygA%2FGlJe6QjYf72oLd%2FADbHH3r68kHTQcX946%2BVsEDjAV0MyefSu6ysv0r8l4ix%2FqAQj6dBt7oPHb%2FGFNL2xp1M1EzChs6OUQ47qVc2k79rvGNvkQs8aC5%2FvzajwMx5lOuxMn9oJmRPG37RE3uI2POJJf03NiElIZbXBL363wnO8pkz%2BaPrJoGvJvhI2acEFW5t8y5hs146Cz1D0ZBp1Rm%2BoNy%2FHZfP9ivfUi1vBV3nAUNsS78ex03If27zIQxJBcRlbq7ak%2Bxgk%2BsESuyx9VDJo8MMs7zywECZlaKr6E02gHJPbOHlFwHzPoHtYuGxYxOFpw9eNQ6g000GGMMQqchBa23%2BWLbrOSb9hArGNLfnxkq0n%2FRVD8QHnZ8xL2KwqOzjYQUmgdUkgJY40V7jkbU7mqFmIjqJWW6psGyi7Q1X8Hh8oubL2DHvxFR292VTG2x6CNB1mFMJYF8MdWIHMZVVZ8g%2FerRXG%2FXqMz5ftDkYSevb3t%2BBcgoKlRTx6nSc6vzYRu4yJW%2BjRMN0T8EO0GtrhGFgFkrb6Ms9D2HA2%2F23YVnCGaGyLag3T5bmzLJzx82r2eijrvP0R3def8mGVMQuWO17pdwpeepGUFu6VThegt6A8KvPqAWwRMSl9kqEaThZwawiXHHRkSVejwIgRiRrEml9zdsBvWZ%2FVPNAzsSrnUBMCxbjZSI4sX6nX0424o9%2BSlUB%2BoCPJcud%2Be6u4okn2WQvvn47BGAtKac1Hzv4LZE4ZwZa%2FunB61%2B7NWAqZ43ZEmVV%2FUWsH40K2jOgEyeHBoQ9PlecJB46sEPMlDL4nRgefk1QWK7DptICl6dxHRk1hsGM%2FOjbnIW4v9q3ILOuWbZI6YRKV93idOPGffaEEBmDWnyj2psMWvrasTKikmUQz1G2J0JJY3z6OWvL2kGV%2FnljF0TTdH7bofSNxLk0FbWAi6ZWQUaD%2FAP%2B%3D&umidGetStatusVal=255&screenPixel=2560x1440&navlanguage=zh-CN&navUserAgent=Mozilla%2F5.0%20%28Macintosh%3B%20Intel%20Mac%20OS%20X%2010_15_7%29%20AppleWebKit%2F537.36%20%28KHTML%2C%20like%20Gecko%29%20Chrome%2F87.0.4280.141%20Safari%2F537.36&navPlatform=MacIntel&appName=taobao&minipara=1%2C1%2C1&appEntrance=taobao_pc&_csrf_token=1WhXEGVzpZkfPhONEGo3K6&umidToken=71fbc9edf087812e8959453cc22590967f38bdeb&hsiz=1c386caaf67de4d7690e183d5f775c4d&newMini2=true&bizParams=&full_redirect=true&style=miniall&appkey=00000000&from=worldlogin&isMobile=false&lang=zh_CN&returnUrl=http%3A%2F%2Fworld.taobao.com%2F&fromSite=0&bx-ua=140#paur3bbqzzF0rzo22Z3TCtSdvWMkKWzgkCqlMM6vlUV6B1j5G5TbbT8jtnpqJRQFCOSTi2n4HZSaCE3n5gy3d8gqlbzx6IUpctgqzzr141HaU6OzzPzbVXlqlbrdNwQ1txyWEAD2Q282UpszaIziVXEFLrfxh8wKJp8Wzl2V8OYNlKMocD+bVJMoLhGP7SvZrI7ZbixvdsbhJ3+tNWB0D7r5Q/6lVY/dLqurkascPcmmze9R9ZssgbxjzFchXZyGe2gFBu4J0WLjXpfr8CbdUdQJXfWMLKY55RVXl/UeJq6p2DA22WdGYDyGKSQJjsGZuCB+Uuyap8ZHqroF9tHFdvtcKqkHBdtasn/5bxx1UryAJRcF2soPguernSVljOA2sEMjNYLgunkj0PoTcy4It0B7NmQsj/4eudPbhIcur0WDqxLv20MtMo59iZexGU3CcXa0m5/+AG5KrjDNUdCoOcgpxKUWuXJfAh6csfSQ3LJn3EADNn9U/VUm56YibY29LAQgKzF9eCOYp39ViX5OThNZiK8prYfoE7nf5aGz3tdS3Bi3mi3S+RZBP6XBZuAxMTBTaQWgZ6uUQ47qVc2k79rvGNvkQs8aC5/vzajwMx5lOuxMnjqnzZPDilBOYyL/NVMGxFLKAj1cuYv7mW8on6gr+Hv+Yj+VFI2juIx1uv4DUX5S2rqnmkPTiD8361zJZw13J55vMEzOxenMbbtkyCVWz4yXsPT7e+Q2B7xdIjGh9RUFYpSd0Ljj/YWVZi87NShkabt2Rf+T2jErB+EpFpSuiQuBnUXQiVTjXY58E0iGaNIm6KbMHX3jfdg2Ia9VYEbEaosrwoCjX/jYgESIH8+jxnppD24cpCgxESTJ7zUupRe1VfuyEwVebVQf4g6Ke82qqCIz+KZrtbE2rnSXURWXgHuwt21S9JAv7wJlqYGfrCbClp46R84caETrBWwqM/YWNcjyeKDsbbEX0lSvZZ1nDAFMdvA5EINAZi/Bw9BrRpil1IfXpTIqXX/4gSCCGtq7rRpUAd0+K40SmcDRz1wOzBf+XH3eHlTDy07a9vrxCzfYS6e7wPD+g9bcryhZaRduvph4SIqSDaieQf+ipV6UaVtmnDPO77X1fgs0yQOVHaupfn7rq8K1DIYm0cXsoNqSshFLbouN0dXDNggYU4b6QYhXbO5lf0qCPyxl7qo0BGev7psZEIvHmTzqVc6YDAfJlIeLqu1Uu9Z6NRw5h+So6dqUCtG4jB8NlCH2IX2RJosoIJLn1TyLYKDpLJwTwcyn7hyP4ilvkhST18Ht4scPNtzR9GhKfSkESqYL64oGtmkXO74Lrr//LakMfY+IiTeOSAVBGgT39vmMzeGrlmVhTadZ3jd0wj212asJ4cR3oR8+bNwntb34qlrTd7Zp7O6hGFZeTwc8jacRBJ6/aznqcF==&bx-umidtoken=T2gAebgRYOzMeGlZ1saou1UjOJaT5_HsJJmNn9LqnShe53XQwi_6wCcC8EiqWka7Cwg=' \\
  --compressed`));
```