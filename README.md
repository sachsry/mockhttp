# Mockhttp
Generic library for mocking http requests in Go

## Purpose
Writing unit tests for http handlers in Go can be a pain. This library is intended to get rid of many of those pain points so you can focus on writing the tests and not on the underlying net/http library. By using generics, this library transforms what can be 20+ lines for a test into 6 or 7 lines. The fluid API design allows for your test definitions to tell the story allowing you to table test your handlers with ease.

```
// Basic test confirming a response payload without the library
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

// Same functionality is getting tested with 5 lines!
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
```

## Examples
See the [examples](https://github.com/sachsry/mockhttp/tree/main/v1/examples) package to see how mockhttp makes writing unit tests for http handlers easy! 

## Features

### Easier syntax for creating requests
Less verbose syntax for creating requests. By using the fluent API, you can easily add path params, contexts, and headers to a request.

IMPORTANT NOTE: Only chi is supported as a mechanism for adding path params, for now.
```
// Basic request without request body
req := mockhttp.NewRequest("GET", "/example", "")

// Set a header that is expected in your handler
req := mockhttp.NewRequest("GET", "/example", "").
         SetHeader("authorization", "bearer 123")

// Set multiple headers expected in your handler
req := mockhttp.NewRequest("GET", "/example", "").
         WithHeaders(map[string]string{
           "authorization": "bearer 123",
           "requestID": "abc-123",
          )}

// Give your request context if needed
req := mockhttp.NewRequest("GET", "/example", "").
          WithValues(map[string]interface{}{
            "tokenID": "123",
          })

// Set chi path params
req := mockhttp.NewRequest("GET", "/things/:id/:name", "").
				WithPathParams(mockhttp.Chi, map[string]string{
					"id":   "1",
					"name": "wax",
				}),
```
### Built in parsing
You don't have to worry about reading the response body via the standard library, or about unmarshaling JSON.
```
// OLD way
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

// Achieve the same result in one line by using generics!
res, err := mockhttp.ToJSONResponse[response.StatusStruct](req.Result())
```

### Handle Errors Same as Regular Responses
Test code looks the same for successful and unsuccessful responses.
```
type MyStruct{
  Age int `json:"age"`
}

type MyError{
  Message string `json:"message"`
}

// Testing a 200
res, err := mockhttp.ToJSONResponse[MyStruct](req.Result())

assert.Nil(t, err)
assert.Equal(t, 200, res.Status())
assert.Equal(t, 23, res.Val.Age) // res.Val is of type MyStruct

// Testing a 400
res, err := mockhttp.ToJSONResponse[MyError](req.Result())

assert.Nil(t, err)
assert.Equal(t, 400, res.Status())
assert.Equal(t, "something bad", res.Val.Message) // res.Val is of type MyError
```

### Create expected responses with a Fluent API
Take advantage of multiple Response types. The `mockhttp.RawResponse` type lets you set an expected status and body. If you are more concerned with getting the expected status code, use this type. Moreover, you can still use the body as you see fit if you want to do string comparisons or otherwise.

If your handlers use JSON, the `mockthttp.JSONResponse` type is a generic type made for testing both success and error responses.
```
// For simple cases
expected := mockhttp.NewRawResponse().
  .WithStatus(200)
  .WithBody("some body")

// If you want to validate JSON, you can specify success or failure
expected := mockhttp.NewJSONResponse[response.StatusStruct]().
  WithSuccess(&response.StatusStruct{Status: "ok"}) // WithSuccess sets status to 200

expected := mockhttp.NewJSONResponse[mockhttp.ServerError]().
		WithFailure(400, &mockhttp.ServerError{
			Status:       "bad request",
			DebugMessage: "something bad",
		})
```
### Table test your API
In the (simple)[https://github.com/sachsry/mockhttp/blob/main/v1/examples/simple_test.go] example, see how the API makes for easy table testing.
```
func TestSimpleHandler(t *testing.T) {
	tests := []mockhttp.TestStruct{ // Use the built in TestStruct for simpler testing scenarios
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

// This handler sends a 200 if the length of the path is even, and a 500 if not.
// For example, for a request mockhttp.NewRequest("GET", "/hello", ""), "/hello" is the path
func handleSimple(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path)%2 != 0 {
		response.Error(w, 500, "length of the path is odd", nil)
		return
	}
	response.Success(w)
}
```

### ADVANCED: Use validation functions for easy, reusable response validation

### ADVANCED: Create your own custom test struct to override service implementationx  