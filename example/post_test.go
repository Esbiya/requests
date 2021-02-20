package example

import (
	"github.com/Esbiya/requests"
	"io/ioutil"
	"log"
	"testing"
)

func TestPostForm(t *testing.T) {
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
}

func TestAsyncPostForm(t *testing.T) {
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
}

func TestSessionPostForm(t *testing.T) {
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
	session := requests.NewSession()
	resp := session.Post(url, headers, data, cookies, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true})
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
}

func TestSessionAsyncPostForm(t *testing.T) {
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
	session := requests.NewSession()
	session.AsyncPost(url, headers, data, cookies, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true}).Then(func(r *requests.Response) {
		if r.Error() != nil {
			log.Fatal(r.Error())
		}
		log.Println(r.Text)
	})
	requests.AsyncWait()
}

func TestPostPayload(t *testing.T) {
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
	resp := requests.Post(url, headers, data, requests.Arguments{Proxy: "http://127.0.0.1:8888", SkipVerifyTLS: true, DisableCompression: false})
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
}

func TestPostBinary(t *testing.T) {
	api := "http://192.168.100.107:7788"
	data, err := ioutil.ReadFile("/Users/esbiya/Desktop/pythonProjects/gitlab/whatsapp/captcha.png")
	if err != nil {
		log.Fatal(err)
	}
	resp := requests.Post(api, data)
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
}

func TestPostFile(t *testing.T) {
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
		Name:  "test.png",                                                        // 文件名
		Param: "fileContent",                                                     // 文件参数名
		Path:  "/Users/esbiya/Desktop/pythonProjects/gitlab/whatsapp/test-0.png", // 文件路径
		Src:   []byte{},                                                          // 文件内容
		Args:  map[string]string{},                                               // 其他参数
	}
	// 支持多表单上传
	files := []*requests.File{file, file1, file2}
	resp := session.Post("https://files-wpa.chat.zalo.me/api/message/upthumb", headers, params, file, file1, file2, files)
	if resp.Error() != nil {
		log.Fatal(resp.Error())
	}
	log.Println(resp.Text)
	log.Println(resp.Cost().String())
}

func TestUpload(t *testing.T) {
	resp, err := requests.Post("http://192.168.0.42:30161/out_upload", &requests.File{
		Path:     "/Users/esbiya/Desktop/javaProjects/app/apps/postern.apk",
		Name:     "postern.apk",
		Param:    "file",
		MimeType: "application/vnd.android.package-archive",
	}).JSON()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
