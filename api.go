package rtu_api

import (
	"fmt"

	"github.com/mchughanalytics/rtu_api/common"
)

func main() {
	fmt.Print("Hello world")
	v, err := common.GetAllVersions()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("versions: %s", v)
}
