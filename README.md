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

​		或者:

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

### session post 请求

```go
resp := session.Post("http://www.baidu.com/", requests.RequestArgs{})
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

