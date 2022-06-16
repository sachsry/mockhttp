package examples_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/sachsry/mockhttp/v1/response"
	"github.com/stretchr/testify/assert"
)

type contextStruct struct {
	ID   int    `json:"id"`
	City string `json:"city"`
}

// Define my own test struct so I can perform validations
// on both success and failure responses from my api
type myTestStruct[S, E any] struct {
	Name    string
	Input   *mockhttp.Request
	Success *mockhttp.JSONResponse[S]
	Error   *mockhttp.JSONResponse[E]
}

func TestErrorsWithContext(t *testing.T) {
	tests := []myTestStruct[contextStruct, mockhttp.ServerError]{
		{
			Name:  "no_id_in_context",
			Input: mockhttp.NewRequest("GET", "/", ""),
			Error: mockhttp.NewJSONResponse[mockhttp.ServerError]().
				WithFailure(400, &mockhttp.ServerError{
					Status:       "bad request",
					DebugMessage: "expected an id of type int in context",
				}).
				// You can use a predefined validation function, or... (see line 63)
				WithValidationFunc(mockhttp.ValidateErrors),
		},
		{
			Name: "no_city_in_context",
			Input: mockhttp.NewRequest("GET", "/", "").
				// Add values to a request's context using this function
				WithValues(map[string]interface{}{
					"id": 123,
				}),
			Error: mockhttp.NewJSONResponse[mockhttp.ServerError]().
				WithFailure(400, &mockhttp.ServerError{
					Status:       "bad request",
					DebugMessage: "expected a city of type string in context",
				}).
				WithValidationFunc(mockhttp.ValidateErrors),
		},
		{
			Name: "success",
			Input: mockhttp.NewRequest("GET", "/", "").
				WithValues(map[string]interface{}{
					"id":   123,
					"city": "Dallas",
				}),
			Success: mockhttp.NewJSONResponse[contextStruct]().
				WithSuccess(&contextStruct{
					ID:   123,
					City: "Dallas",
				}).
				// You can define in line validations for tests
				WithValidationFunc(func(expected, result contextStruct) error {
					if expected.ID != result.ID || expected.City != result.City {
						return errors.New("unexpected result")
					}
					return nil
				}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			handleRequestWithContext(tt.Input.W, tt.Input.R)

			// handle error test cases
			if tt.Error != nil {
				res, err := mockhttp.ToJSONResponse[mockhttp.ServerError](tt.Input.Result())
				assert.Nil(t, err)
				assert.NotNil(t, res)
				tt.Error.Validate(res)
			}

			// handle success test cases
			if tt.Success != nil {
				res, err := mockhttp.ToJSONResponse[contextStruct](tt.Input.Result())
				assert.Nil(t, err)
				assert.NotNil(t, res)
				tt.Success.Validate(res)
			}
		})
	}
}

// This handler simply looks for values in the context of certain types
// and errors out if it doesn't find either
func handleRequestWithContext(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(int)
	if !ok {
		response.Error(w, 400, "expected an id of type int in context", nil)
		return
	}
	city, ok := r.Context().Value("city").(string)
	if !ok {
		response.Error(w, 400, "expected a city of type string in context", nil)
		return
	}
	vals := map[string]interface{}{
		"id":   id,
		"city": city,
	}
	response.SuccessWithBody(w, vals)
}
