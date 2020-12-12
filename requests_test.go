package requests

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	resp := Get("https://www.baidu.com", RequestArgs{}).Text
	log.Println(resp)
}
