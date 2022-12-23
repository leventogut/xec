package xec

// HandleError receives errors, stores and returns them.
func HandleError(err error) {
	if err != nil {
		err
	}
	return
}
