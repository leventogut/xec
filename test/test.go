package main

import (
	"fmt"
	"os"
)

func main() {

	var envKeyValue string
	var envKeyValues []string
	envKeyValues = os.Environ()

	for _, envKeyValue = range envKeyValues {
		fmt.Printf("%v\n", envKeyValue)
	}
}
