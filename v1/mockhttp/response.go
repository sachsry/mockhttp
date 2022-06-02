package mockhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response interface {
	Status() int
	Body() string
}

type RawResponse struct {
	status int
	body   string
}

func NewRawResponse() *RawResponse {
	return &RawResponse{}
}

func (r *RawResponse) Status() int {
	return r.status
}

func (r *RawResponse) Body() string {
	return r.body
}

func (r *RawResponse) WithStatus(status int) *RawResponse {
	r.status = status
	return r
}

type JSONResponse[T any] struct {
	status         int
	body           string
	Val            *T
	validationFunc func(expected, result T) error
}

func NewJSONResponse[T any]() *JSONResponse[T] {
	return &JSONResponse[T]{}
}

func (r *JSONResponse[T]) Status() int {
	return r.status
}

func (r *JSONResponse[T]) Body() string {
	return r.body
}

func (r *JSONResponse[T]) WithStatus(status int) *JSONResponse[T] {
	r.status = status
	return r
}

func (r *JSONResponse[T]) WithSuccess(val *T) *JSONResponse[T] {
	r.status = 200
	r.Val = val
	return r
}

func (r *JSONResponse[T]) WithFailure(status int, val *T) *JSONResponse[T] {
	r.status = status
	r.Val = val
	return r
}

func (r *JSONResponse[T]) WithValidationFunc(f func(expected, result T) error) *JSONResponse[T] {
	r.validationFunc = f
	return r
}

// ToResponse takes a httpResponse and maps it to a JSONResponse object
// It saves the status and saves the body
func ToResponse(res *http.Response) (*RawResponse, error) {
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &RawResponse{
		status: res.StatusCode,
		body:   string(data),
	}, nil
}

// ToJSONResponse takes a httpResponse and maps it to a JSONResponse object
// It saves the status and parses the body into the expected type
// It will return an error if the http.Response is not a 200
func ToJSONResponse[T any](res *http.Response) (*JSONResponse[T], error) {
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ret := &JSONResponse[T]{
		status: res.StatusCode,
		body:   string(data),
	}

	if len(ret.body) == 0 {
		return nil, errors.New("expected a payload in the response body, but got an empty string")
	}

	var t T
	ret.Val = &t
	err = json.Unmarshal(data, ret.Val)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Validate performs validation for two JSON Responses
func (expected *JSONResponse[T]) Validate(result *JSONResponse[T]) error {
	if expected == nil {
		return errors.New("receiver expected should not be nil")
	}
	if result == nil {
		return errors.New("parameter result should not be nil")
	}
	if expected.status != result.status {
		return fmt.Errorf("expected status %d, but got %d", expected.status, result.status)
	}
	if expected.validationFunc != nil {
		return expected.validationFunc(*expected.Val, *result.Val)
	}
	return nil
}
