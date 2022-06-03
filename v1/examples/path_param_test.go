package examples_test

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/sachsry/mockhttp/v1/response"
	"github.com/stretchr/testify/assert"
)

func TestChiPathParams(t *testing.T) {
	tests := []mockhttp.TestStruct{
		{
			Name:     "no_params",
			Input:    mockhttp.NewRequest("GET", "/things/:id", ""),
			Expected: mockhttp.NewRawResponse().WithStatus(500),
		},
		{
			Name: "id_param",
			Input: mockhttp.NewRequest("GET", "/things/:id", "").
				WithPathParams(mockhttp.Chi, map[string]string{
					"id": "1",
				}),
			Expected: mockhttp.NewRawResponse().
				WithStatus(200).
				WithBody(`{"id":"1"}`),
		},
		{
			Name: "multiple_params",
			Input: mockhttp.NewRequest("GET", "/things/:id/:name", "").
				WithPathParams(mockhttp.Chi, map[string]string{
					"id":   "1",
					"name": "wax",
				}),
			Expected: mockhttp.NewRawResponse().
				WithStatus(200).
				WithBody(`{"id":"1","name":"wax"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			handleChiPathParams(tt.Input.W, tt.Input.R)

			res, err := mockhttp.ToResponse(tt.Input.Result())

			assert.Nil(t, err)
			assert.Equal(t, tt.Expected.Status(), res.Status())
			if tt.Expected.Body() != "" {
				assert.Equal(t, tt.Expected.Body(), res.Body())
			}
		})
	}
}

func handleChiPathParams(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		response.Error(w, 500, "no id provided", nil)
		return
	}
	ret := map[string]string{
		"id": id,
	}
	name := chi.URLParam(r, "name")
	if name != "" {
		ret["name"] = name
	}
	response.SuccessWithBody(w, ret)
}
