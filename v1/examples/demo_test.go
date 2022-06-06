package examples_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/sachsry/mockhttp/v1/response"
	"github.com/stretchr/testify/assert"
)

func TestWithoutLib_Basic(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	w := httptest.NewRecorder()
	handleSuccess(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.Equal(t, `{"status":"ok"}`, string(data))
}

func TestWithoutLib_Parsing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	w := httptest.NewRecorder()
	handleSuccess(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response: %v", err)
	}

	var ok response.StatusStruct
	err = json.Unmarshal(data, &ok)
	if err != nil {
		t.Fatalf("error unmarshalling data: %v", err)
	}

	assert.Equal(t, "ok", ok.Status)
}

func TestWithLib_Basic(t *testing.T) {
	req := mockhttp.NewRequest("GET", "/example", "")
	handleSuccess(req.W, req.R)

	res, err := mockhttp.ToResponse(req.Result())
	assert.Nil(t, err)
	assert.Equal(t, `{"status":"ok"}`, res.Body())
}

func TestWithLib_Parsing(t *testing.T) {
	req := mockhttp.NewRequest("GET", "/example", "")
	handleSuccess(req.W, req.R)

	res, err := mockhttp.ToJSONResponse[response.StatusStruct](req.Result())
	assert.Nil(t, err)
	assert.Equal(t, "ok", res.Val.Status)
}

func handleSuccess(w http.ResponseWriter, r *http.Request) {
	// response.Success writes a 200 and sends a {status: ok} json response
	response.Success(w)
}
