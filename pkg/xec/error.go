package xec

import "fmt"

// HandleError receives errors, stores and returns them.
func HandleError(err error) {
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
	return
}
