package mockhttp_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/sachsry/mockhttp/v1/response"
	"github.com/stretchr/testify/assert"
)

func TestRawResponse_New(t *testing.T) {
	res := mockhttp.NewRawResponse()
	assert.NotNil(t, res)

	res = res.WithStatus(300)
	assert.Equal(t, 300, res.Status())
	assert.Equal(t, "", res.Body())
}

func TestRawResponse_SuccessHandler(t *testing.T) {
	httpReq := mockhttp.NewRequest("GET", "/", "")
	successHandler(httpReq.W, httpReq.R)

	res, err := mockhttp.ToResponse(httpReq.Result())
	if err != nil {
		t.Fatalf("expected error to be nil, but found an error: %v", err.Error())
	}

	assert.Equal(t, 200, res.Status())
	assert.Equal(t, `{"status":"ok"}`, res.Body())
}

func TestRawResponse_FailureHandler(t *testing.T) {
	httpReq := mockhttp.NewRequest("GET", "/", "")
	failHandler(httpReq.W, httpReq.R)

	res, err := mockhttp.ToResponse(httpReq.Result())
	if err != nil {
		t.Fatalf("expected error to be nil, but found an error: %v", err.Error())
	}

	assert.Equal(t, 400, res.Status())
	assert.Equal(t, `{"message":"something bad","status":"bad request"}`, res.Body())
}

func TestJSONResponse_New(t *testing.T) {
	res := mockhttp.NewJSONResponse[response.StatusStruct]()
	res = res.WithSuccess(&response.StatusStruct{Status: "ok"})

	assert.Equal(t, 200, res.Status())
	assert.Equal(t, "ok", res.Val.Status)

	res = res.WithStatus(222)
	assert.Equal(t, 222, res.Status())
}

func TestJSONResponse_SuccessHandler(t *testing.T) {
	httpReq := mockhttp.NewRequest("GET", "/", "")
	successHandler(httpReq.W, httpReq.R)

	res, err := mockhttp.ToJSONResponse[response.StatusStruct](httpReq.Result())

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, `{"status":"ok"}`, res.Body())
	assert.Equal(t, "ok", res.Val.Status)
}

func TestJSONResponse_FailureHandler(t *testing.T) {
	httpReq := mockhttp.NewRequest("GET", "/", "")
	failHandler(httpReq.W, httpReq.R)

	res, err := mockhttp.ToJSONResponse[mockhttp.ServerError](httpReq.Result())

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, `{"message":"something bad","status":"bad request"}`, res.Body())
	assert.Equal(t, "something bad", res.Val.DebugMessage)
	assert.Equal(t, "bad request", res.Val.Status)
}

func TestJSONResponse_NoBodyError(t *testing.T) {
	httpReq := mockhttp.NewRequest("GET", "/", "")
	nothingHandler(httpReq.W, httpReq.R)

	res, err := mockhttp.ToJSONResponse[mockhttp.ServerError](httpReq.Result())

	assert.Nil(t, res)
	assert.Equal(t, "expected a payload in the response body, but got an empty string", err.Error())
}

func TestJSONResponse_UnmarshalError(t *testing.T) {
	httpReq := mockhttp.NewRequest("GET", "/", "")
	failHandler(httpReq.W, httpReq.R)

	res, err := mockhttp.ToJSONResponse[struct {
		Status int `json:"status"`
	}](httpReq.Result())

	assert.Nil(t, res)
	assert.Equal(t, "json: cannot unmarshal string into Go struct field .status of type int", err.Error())
}

func TestJSONResponse_ValidationFunc_Success(t *testing.T) {
	expected := mockhttp.NewJSONResponse[response.StatusStruct]().
		WithSuccess(&response.StatusStruct{Status: "ok"}).
		WithValidationFunc(func(expected, result response.StatusStruct) error {
			if expected.Status != result.Status {
				return fmt.Errorf("unexpected status found in result: %v", result.Status)
			}
			return nil
		})

	httpReq := mockhttp.NewRequest("GET", "/", "")
	successHandler(httpReq.W, httpReq.R)
	result, err := mockhttp.ToJSONResponse[response.StatusStruct](httpReq.Result())

	assert.Nil(t, err)
	assert.NotNil(t, result)

	validationErr := expected.Validate(result)
	assert.Nil(t, validationErr)
}

func TestJSONResponse_ValidationFunc_ServerError(t *testing.T) {
	expected := mockhttp.NewJSONResponse[mockhttp.ServerError]().
		WithFailure(400, &mockhttp.ServerError{
			Status:       "bad request",
			DebugMessage: "something bad",
		}).
		WithValidationFunc(mockhttp.ValidateErrors)

	httpReq := mockhttp.NewRequest("GET", "/", "")
	failHandler(httpReq.W, httpReq.R)
	result, err := mockhttp.ToJSONResponse[mockhttp.ServerError](httpReq.Result())

	assert.Nil(t, err)
	assert.NotNil(t, result)

	validationErr := expected.Validate(result)
	assert.Nil(t, validationErr)
}

func TestValidate_NullExpected(t *testing.T) {
	var expected *mockhttp.JSONResponse[any]

	err := expected.Validate(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "receiver expected should not be nil", err.Error())
}

func TestValidate_NullResult(t *testing.T) {
	var expected mockhttp.JSONResponse[any]

	err := expected.Validate(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "parameter result should not be nil", err.Error())
}

func TestValidate_Error(t *testing.T) {
	expected := mockhttp.NewJSONResponse[response.StatusStruct]().
		WithSuccess(&response.StatusStruct{Status: "ok"}).
		WithValidationFunc(func(expected, result response.StatusStruct) error {
			if expected.Status != result.Status {
				return fmt.Errorf("unexpected status found in result: %v", result.Status)
			}
			return nil
		})
	result := mockhttp.NewJSONResponse[response.StatusStruct]().
		WithSuccess(&response.StatusStruct{Status: "not ok"})

	err := expected.Validate(result)

	assert.NotNil(t, err)
	assert.Equal(t, "unexpected status found in result: not ok", err.Error())
}

func TestValidate_NoValidationFunc(t *testing.T) {
	expected := mockhttp.NewJSONResponse[response.StatusStruct]().
		WithSuccess(&response.StatusStruct{Status: "ok"})
	result := mockhttp.NewJSONResponse[response.StatusStruct]().
		WithSuccess(&response.StatusStruct{Status: "not ok"})

	err := expected.Validate(result)

	assert.Nil(t, err)
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	response.Success(w)
}

func failHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, 400, "something bad", nil)
}

func nothingHandler(w http.ResponseWriter, r *http.Request) {}
