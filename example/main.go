package main

import (
	"github.com/mengxiaozhu/aligw"
	"os"
	"log"
)

func main() {

	gw := &aligw.AliGateway{AppKey: os.Getenv("appKey"), AppSecret: os.Getenv("appSecret")}
	str, err := gw.Get(os.Getenv("host"), os.Getenv("path")).Send().String()
	log.Println(err, str)
}
