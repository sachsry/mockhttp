package mockhttp

import "testing"

type TestStruct[T any] struct {
	Name           string
	Input          *Request
	Expected       *Response[T]
	ValidationFunc func(t *testing.T, expected, result *T)
}
