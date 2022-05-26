package examples_test

import (
	"testing"

	"github.com/sachsry/mockhttp/v1/examples"
	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/sachsry/mockhttp/v1/response"
	"github.com/stretchr/testify/assert"
)

func TestSimpleHandler(t *testing.T) {
	tests := []mockhttp.TestStruct[response.StatusStruct]{
		{
			Name:  "even length",
			Input: mockhttp.NewRequest("GET", "/suh", ""), // This is even because the "/" counts as a char in the path
			Expected: &mockhttp.Response[response.StatusStruct]{
				Status: 200,
				Struct: &response.StatusStruct{Status: "ok"},
			},
		},
		{
			Name:  "odd length",
			Input: mockhttp.NewRequest("GET", "/blah", ""), // This is odd because the "/" counts as a char in the path
			Expected: &mockhttp.Response[response.StatusStruct]{
				Status: 500,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			examples.HandleSimple(tt.Input.W, tt.Input.R)

			res, err := mockhttp.ToResponse[response.StatusStruct](tt.Input.W.Result())

			assert.Nil(t, err, res)

			// Define custom validation function for validating a 200 response
			validationFunc := func(t *testing.T, expected, result *response.StatusStruct) {
				assert.NotNil(t, expected)
				assert.NotNil(t, result)
				assert.Equal(t, expected.Status, result.Status)
			}

			tt.Expected.Validate(t, res, validationFunc)

			// Perform validations on error responses
			if tt.Expected.Status >= 400 {
				var se mockhttp.ServerError
				err := res.UnmarshalError(&se)
				if err != nil {
					t.Errorf("expected unmarshal error to be nil, instead got %v", err)
				}
				assert.Equal(t, "internal error", se.Status)
				assert.Equal(t, "length of the path is odd", se.DebugMessage)
			}
		})
	}
}
