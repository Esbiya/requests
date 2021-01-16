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
  	resp := requests.Get(url, requests.RequestArgs{})
	  if resp.StatusCode != http.StatusOK {
		    log.Fatal("状态码异常")
   	}
   	log.Println(resp.Text)
}
```

### get 请求

```go
url := "https://www.baidu.com/"
params := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
}
resp := requests.Get(url, requests.RequestArgs{
	Headers: headers,
	Params: Params,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 异步 get 请求

```
c := make(*Response, 1)
c := make(chan *Response, 1)
requests.AsyncGet(url, RequestArgs{}, c)
log.Println((<-c).JSON())
```

### post form 请求

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	Data: data,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 异步 post 请求

```
c := make(*Response, 1)
c := make(chan *Response, 1)
requests.AsyncPost(url, RequestArgs{}, c)
log.Println((<-c).JSON())
```

### post json 请求

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	JSON: data,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 请求头设置

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx", 
    "Connection": "keep-alive",
    "Cookie": "111=222; 333=444",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	JSON: data,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### cookie 设置

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
	"Cookie": "111=222; 333=444",  // 设置方式 1
}
Cookies := requests.DataMap{
	"111": "222", 
    "333": "444",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	Data: data,
	Cookies: cookies,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 代理设置

- 支持 s5/http/https

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	Data: data,
	Proxy: "http://127.0.0.1:8888",
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 禁止重定向

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	Data: data,
	AllowRedirects: false,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 禁止响应数据压缩

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	Data: data,
	DisableCompression: true,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 禁止 tcp 链接复用

```go
url := "https://www.baidu.com/"
data := requests.DataMap{
	"hello": "there",
}
headers := requests.DataMap{
	"User-Agent": "xxx",
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	Data: data,
    DisableKeepAlive: true,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
```

### 上传文件

```go
url := "https://www.baidu.com/"
headers := requests.DataMap{
	"User-Agent": "xxx",
}
filaName := "xxx.jpg"
b, _ := ioutil.ReadFile(fileName)
file := requests.File(fileName, b)
fileData := requests.DataMap{
	"file": file,
	"Content-Type": "image/jpeg"
}
resp := requests.Post(url, requests.RequestArgs{
	Headers: headers,
	File: fileData,
})
if resp.StatusCode != http.StatusOK {
	log.Fatal("状态码异常")
}
log.Println(resp.Text)
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
resp := session.Get("http://www.baidu.com/", requests.RequestsArgs{})
```

### 异步 session get 请求

```
c := make(*Response, 1)
c := make(chan *Response, 1)
session.AsyncGet(url, RequestArgs{}, c)
log.Println((<-c).JSON())
```

### session post 请求

```go
resp := session.Post("http://www.baidu.com/", requests.RequestArgs{})
```

### 异步 session post 请求

```
c := make(*Response, 1)
c := make(chan *Response, 1)
session.AsyncPost(url, RequestArgs{}, c)
log.Println((<-c).JSON())
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

```
url := session.GetUrl()
proxy := session.GetProxy()
timeout := session.GetTimeout()
ssl := session.GetSkipVerifyTLS()
allowRedirects := session.GetAllowRedirects()
disableKeepAlive := session.GetDisableKeepAlive()
disableCompression := session.GetDisableCompression()
```

### session 复制

```
session1 := session.Copy()
```

### session 中间件注册

* 请求前对请求参数预处理, 如对请求参数进行排序加签操作

```
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
```

* 请求前对请求对象进行预处理

```
_ = session.RegisterBeforeReqHook(func(request *http.Request) error {
	request.Close = true
	return nil
})
```

* 请求完成对响应进行预处理, 如将加密响应解密成明文

```
_ = session.RegisterAfterRespHook(func(response *Response) error {
	var err error
	response.Bytes, err = base64.StdEncoding.DecodeString(response.Text)
	response.Text = string(response.Bytes)
	return err
})
```

### 响应

```go
resp := session.Post("http://www.baidu.com/", requests.RequestArgs{})
resp.StatusCode          // 状态码
resp.Bytes               // 响应字节
resp.Text                // 响应字符串
resp.JSON()              // 解析 json 响应
resp.CallbackJSON()      // 解析 jQuery({"1": "2"}) json 响应
resp.SetEncode("gbk")    // 设置响应编码
resp.GetEncode()         // 获取响应编码
resp.SaveFile("xxx.jpg") // 文件写入
resp.Err                 // 响应错误, 正常响应为 nil
resp.Header              // 响应头
resp.Cookies()           // 响应 cookies, []*http.Cookie
resp.Location()          // 跳转 url
resp.ContentLength       // 响应内容大小
resp.Close               // tcp 连接是否已关闭 bool
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