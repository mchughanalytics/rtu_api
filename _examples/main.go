package main

import (
	"log"

	"github.com/mchughanalytics/rtu_api"
)

func main() {

	key := rtu_api.NewGUID()
	rtu_api.OpenBrowser("index.html")
	// cli, err := rtu_api.NewRmClient("10.11.99.1", "root", "9Cdoe1dZQs")
	// if err != nil {
	// 	panic(err)
	// }
	log.Println(key)

	rtu := rtu_api.NewRtuState()
	rtu.Init(0, true, key)

	log.Println("Closing.")

}
