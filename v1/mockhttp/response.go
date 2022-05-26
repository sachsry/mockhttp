package mockhttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Response represents a json http response
type Response[T any] struct {
	Status int
	Body   string
	Header http.Header
	Struct *T
}

// ToResponse takes a httpResponse and maps it to a Response object
// It saves the status and parses the body into either the expected 200 response
// or into the server error 400+ response
func ToResponse[T any](res *http.Response) (*Response[T], error) {
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ret := &Response[T]{
		Status: res.StatusCode,
		Body:   string(data),
		Header: res.Header,
	}

	// Need to check the length here since some handlers don't write responses to W
	if ret.Status == 200 && len(string(data)) > 0 {
		var t T
		ret.Struct = &t
		err := json.Unmarshal(data, ret.Struct)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

// Unmarshal error converts the response body into an error type of your choice
// Supply a pointer to the struct you would like to unmarshal the error into
func (r *Response[T]) UnmarshalError(out interface{}) error {
	if r.Status > 399 {
		return json.Unmarshal([]byte(r.Body), out)
	}
	return nil
}

// Validate is a generic validation function to be used in testing.
// Expected is a receiver that is the expected response from the handler
// Result is the result the handler returns in the test.
// Finally, a validation func is given to perform validations on the resulting struct, T
func (expected *Response[T]) Validate(t *testing.T, result *Response[T], validationFunc func(t *testing.T, expected, result *T)) {
	assert.Equal(t, expected.Status, result.Status)
	if expected.Status == 200 {
		validationFunc(t, expected.Struct, result.Struct)
		return
	}
}
