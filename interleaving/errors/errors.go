package errors

import "errors"

var (
	NotSliceError        = errors.New("Passed ranking is not slice")
	NotImplementKeyError = errors.New("Please implement Key")
)
