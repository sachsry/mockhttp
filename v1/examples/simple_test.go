package examples_test

import (
	"net/http"
	"testing"

	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/sachsry/mockhttp/v1/response"
	"github.com/stretchr/testify/assert"
)

func TestSimpleHandler(t *testing.T) {
	tests := []mockhttp.TestStruct{
		{
			Name:     "even length",
			Input:    mockhttp.NewRequest("GET", "/suh", ""), // This is even because the "/" counts as a char in the path
			Expected: mockhttp.NewRawResponse().WithStatus(200),
		},
		{
			Name:     "odd length",
			Input:    mockhttp.NewRequest("GET", "/blah", ""), // This is odd because the "/" counts as a char in the path
			Expected: mockhttp.NewRawResponse().WithStatus(500),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			handleSimple(tt.Input.W, tt.Input.R)

			res, err := mockhttp.ToResponse(tt.Input.Result())
			assert.Nil(t, err, res)
			assert.Equal(t, tt.Expected.Status(), res.Status())
		})
	}
}

func handleSimple(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path)%2 != 0 {
		response.Error(w, 500, "length of the path is odd", nil)
		return
	}
	response.Success(w)
}
