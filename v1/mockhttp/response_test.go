package mockhttp_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/stretchr/testify/assert"
)

func TestHandleEmptySuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}

	r := mockhttp.NewRequest("GET", "/example", "")
	handler(r.W, r.R)

	result, err := mockhttp.ToResponse[any](r.W.Result())

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Status, 200)
}

type temp struct {
	I   int    `json:"i"`
	Str string `json:"str"`
}

func TestHandleSuccessWithStruct(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		s := temp{I: 14, Str: "suh"}
		bs, err := json.Marshal(s)
		if err != nil {
			panic("unexpected error occurred during test")
		}
		w.WriteHeader(200)
		w.Write(bs)
	}

	r := mockhttp.NewRequest("GET", "/example", "")
	handler(r.W, r.R)

	result, err := mockhttp.ToResponse[temp](r.W.Result())

	assert.Nil(t, err)
	assert.Equal(t, result.Status, 200)
	assert.Equal(t, 14, result.Struct.I)
	assert.Equal(t, "suh", result.Struct.Str)
}

func TestHandleSuccessWithMap(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		s := map[string]string{
			"a": "alpha",
			"b": "beta",
		}
		bs, err := json.Marshal(s)
		if err != nil {
			panic("unexpected error occurred during test")
		}
		w.WriteHeader(200)
		w.Write(bs)
	}

	r := mockhttp.NewRequest("GET", "/example", "")
	handler(r.W, r.R)

	result, err := mockhttp.ToResponse[map[string]string](r.W.Result())

	assert.Nil(t, err)
	assert.Equal(t, result.Status, 200)
	assert.NotNil(t, result.Struct)

	rmap := *result.Struct
	assert.Equal(t, "alpha", rmap["a"])
	assert.Equal(t, "beta", rmap["b"])
}

func TestHandleSuccessWithSlice(t *testing.T) {
	s := []temp{
		{I: 1, Str: "one"},
		{I: 2, Str: "two"},
		{I: 3, Str: "three"},
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		bs, err := json.Marshal(s)
		if err != nil {
			panic("unexpected error occurred during test")
		}
		w.WriteHeader(200)
		w.Write(bs)
	}

	r := mockhttp.NewRequest("GET", "/example", "")
	handler(r.W, r.R)

	result, err := mockhttp.ToResponse[[]temp](r.W.Result())

	assert.Nil(t, err)
	assert.Equal(t, result.Status, 200)
	assert.NotNil(t, result.Struct)

	rstruct := *result.Struct
	assert.Equal(t, 3, len(rstruct))
	for i := 0; i < len(s); i++ {
		assert.Equal(t, s[i].I, rstruct[i].I)
		assert.Equal(t, s[i].Str, rstruct[i].Str)
	}
}

func TestUnmarshalError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		s := temp{I: 14, Str: "suh"}
		bs, err := json.Marshal(s)
		if err != nil {
			panic("unexpected error occurred during test")
		}
		w.WriteHeader(200)
		w.Write(bs)
	}

	r := mockhttp.NewRequest("GET", "/example", "")
	handler(r.W, r.R)

	result, err := mockhttp.ToResponse[[]temp](r.W.Result())

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "json: cannot unmarshal object into Go value of type []mockhttp_test.temp")
}

func TestHandleEmptyFailure(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}

	r := mockhttp.NewRequest("GET", "/example", "")
	handler(r.W, r.R)

	result, err := mockhttp.ToResponse[any](r.W.Result())

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Status, 500)
	assert.Nil(t, result.Struct)
}

func TestHandleFailureWithStruct(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		se := mockhttp.ServerError{
			Status:       "internal error",
			DebugMessage: "something went wrong",
		}
		bs, err := json.Marshal(se)
		if err != nil {
			panic("unexpected error occurred in test")
		}
		w.Write(bs)
	}

	r := mockhttp.NewRequest("GET", "/example", "")
	handler(r.W, r.R)

	result, err := mockhttp.ToResponse[any](r.W.Result())

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Status, 500)
	assert.Nil(t, result.Struct)

	var se mockhttp.ServerError
	err = result.UnmarshalError(&se)
	assert.Nil(t, err)
	assert.Equal(t, "internal error", se.Status)
	assert.Equal(t, "something went wrong", se.DebugMessage)
}
