package requests

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	session := NewSession().SetProxy("http://127.0.0.1:8888").SetSkipVerifyTLS(true)

	sidCookies, _ := TransferCookies([]map[string]interface{}{
		{
			"Domain": "zalo.me", "HttpOnly": true, "MaxAge": 31535990,
			"Name": "zpw-sek", "Path": "/",
			"Raw":        "zpw-sek=9O66.306388517.a0.aNrswA0G_KQ53hvKWH3zHjioc2s2AjyrxaUH8v1-gpxc8QuyvKQd3l8bqdR7BS1atIqJAQmDsux54G3RNayaE0; Path=/; Max-Age=31535990; Domain=zalo.me; HttpOnly; Secure;",
			"RawExpires": "", "SameSite": 4, "Secure": true, "Unparsed": nil,
			"Value": "9O66.306388517.a0.aNrswA0G_KQ53hvKWH3zHjioc2s2AjyrxaUH8v1-gpxc8QuyvKQd3l8bqdR7BS1atIqJAQmDsux54G3RNayaE0",
		}, {
			"Domain": "zalo.me", "HttpOnly": true, "MaxAge": 31535990,
			"Name": "zpw-sekm", "Path": "/",
			"Raw":        "zpw-sekm=V2aY.306388517.240.lJ-bTyPRcn7ZHTW3orkR4RCe_cFrQwqY_cEk8F5cTNxAHhtPerfcZ3W4cn4; Path=/; Max-Age=31535990; Domain=zalo.me; HttpOnly; Secure;",
			"RawExpires": "", "SameSite": 0, "Secure": true, "Unparsed": nil,
			"Value": "V2aY.306388517.240.lJ-bTyPRcn7ZHTW3orkR4RCe_cFrQwqY_cEk8F5cTNxAHhtPerfcZ3W4cn4",
		},
	})
	session.SetCookies("https://chat.zalo.me/", sidCookies)

	headers := DataMap{
		"Connection": "keep-alive",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
	}
	resp := session.Get("https://p4-msg.chat.zalo.me/?zpw_ver=72&zpw_type=30&params=2002847889307%2C2002847889307%2CeyIxIjp7ImV2aWN0IjoxLCJpZHMiOlsiMTk1ODUwMzM2Njc2OSIsIjE5NzM5ODE0MzcwODEiLCIxOTczOTgxNTMzMDk1IiwiMTk3NDE0NDI0MjY3NSIsIjE5NzQxNDQzMjYyNjciLCIxOTc3Njk1Mjc3ODE1IiwiMTk3Nzk3OTk2NDUzNCIsIjE5NzgzNDU4MzgyMTIiLCIxOTgwNTM5OTc5OTM1IiwiMTk4MDU0NTE4NDgxOSJdLCJsYXN0SWQiOiIxOTgwNTQ1MTg0ODE5In0sIjIiOnsiZXZpY3QiOjEsImlkcyI6WyIxOTk3MDM5MTk4MjI2IiwiMjAwMTk1OTE1NzYyMSIsIjIwMDE5ODkxOTc5MDIiLCIyMDAyMzY1ODk1NjIxIiwiMjAwMjM2Nzk1NDQ3OSIsIjIwMDI1MTIyMDMxMTIiLCIyMDAyNzUxMTk1NTMxIiwiMjAwMjgwNzI2NTg2MyIsIjIwMDI4MzE5MjM4NzQiLCIyMDAyODQ3ODg5MzA3Il0sImxhc3RJZCI6IjIwMDI4NDc4ODkzMDcifSwiMyI6eyJldmljdCI6MSwiaWRzIjpbIjE5NDc3NTA2MTY1NTEiLCIxOTQ3NzUwNjE2NTUyIiwiMTk0Nzc1MDYxNjU1MyIsIjE5NDc4NDg3NTczMDgiLCIxOTQ3ODQ4NzcwMDUxIiwiMTk0Nzg0ODc3MDA1MiIsIjE5NDc4NDg3NzAwNTMiLCIxOTQ3ODQ4NzcwMDU0IiwiMTk0Nzg0ODc3MDA1NSIsIjE5NDc4NDg3NzAxMzYiXSwibGFzdElkIjoiMTk0Nzg0ODc3MDEzNiJ9fQ%3D%3D%2C1&ts=1610508723", RequestArgs{
		Headers: headers,
	})
	if resp.Err != nil {
		log.Fatal(resp.Err)
	}
	log.Println(resp.Text)
}
