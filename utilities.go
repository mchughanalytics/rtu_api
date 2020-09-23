package rtu_api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
)

func ToJSON(i interface{}) []byte {
	out, _ := json.MarshalIndent(i, "", "  ")
	return out
}

func NewGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	guid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return guid
}

func TransformStructs(input, output interface{}) error {

	inputBytes := ToJSON(input)
	err := json.Unmarshal(inputBytes, output)

	return err
}
