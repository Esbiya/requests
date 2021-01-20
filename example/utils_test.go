package example

import (
	"github.com/Esbiya/requests"
	"log"
	"testing"
)

func TestRandomUserAgent(t *testing.T) {
	ua := requests.RandomUserAgent(requests.Chrome)
	log.Println(ua)
	ua1 := requests.RandomUserAgent(requests.Safari)
	log.Println(ua1)
	ua2 := requests.RandomUserAgent(requests.IE)
	log.Println(ua2)
	ua3 := requests.RandomUserAgent(requests.Opera)
	log.Println(ua3)
	ua4 := requests.RandomUserAgent(nil)
	log.Println(ua4)
}
