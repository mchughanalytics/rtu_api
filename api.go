package rtu_api

import (
	"fmt"
	"E:/Source/RTU_API/src/common"
)

func main() {
	fmt.Print("Hello world")
	v, err := common.GetAllVersions()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("versions: %s", v)
}