package main

import (
	"log"

	"github.com/mchughanalytics/rtuapi"
)

func main() {

	key := rtuapi.NewGUID()
	//rtuapi.OpenBrowser("index.html")
	// cli, err := rtuapi.NewRmClient("10.11.99.1", "root", "9Cdoe1dZQs")
	// if err != nil {
	// 	panic(err)
	// }
	log.Println(key)

	rtu := rtuapi.NewRtuState()
	rtu.Init(0, true, key)

	log.Println("Closing.")

}
