package examples_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sachsry/mockhttp/v1/examples/service"
	"github.com/sachsry/mockhttp/v1/examples/service/servicefakes"
	"github.com/sachsry/mockhttp/v1/mockhttp"
	"github.com/sachsry/mockhttp/v1/response"
	"github.com/stretchr/testify/assert"
)

// Here we are simulating an API that has a dependency on a service to read/write to a database
type myAPI struct {
	s service.Service
}

// In practice, this would be the PATCH /profiles/:id handler
func (m *myAPI) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	err := m.s.UpdateProfile("someVal")
	if err != nil {
		response.Error(w, 500, "unable to update profile", err)
		return
	}
	response.Success(w)
}

// you can create a similar struct for your use case
// it is likely you'll have more than one service to mock
type mocks struct {
	fs *servicefakes.FakeService
}

func NewMocks() mocks {
	return mocks{
		fs: &servicefakes.FakeService{},
	}
}

type myMockTestStruct struct {
	Name          string
	Input         *mockhttp.Request
	Expected      mockhttp.Response
	OverrideMocks func(m mocks) // pass a function to the test cases to override an interface
}

func TestMyApi(t *testing.T) {
	tests := []myMockTestStruct{
		{
			Name:     "happy_path",
			Input:    mockhttp.NewRequest("GET", "/", ""),
			Expected: mockhttp.NewRawResponse().WithStatus(200),
			// No overriding needed for the happy path because counterfeiter library has defaults
		},
		{
			Name:     "sad_path",
			Input:    mockhttp.NewRequest("GET", "/", ""),
			Expected: mockhttp.NewRawResponse().WithStatus(500),
			OverrideMocks: func(m mocks) {
				m.fs.UpdateProfileReturns(errors.New("something bad"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mocks := NewMocks()
			myapi := myAPI{s: mocks.fs}
			if tt.OverrideMocks != nil {
				tt.OverrideMocks(mocks)
			}

			myapi.handleProfileUpdate(tt.Input.W, tt.Input.R)

			res, err := mockhttp.ToResponse(tt.Input.Result())
			assert.Nil(t, err)
			assert.Equal(t, tt.Expected.Status(), res.Status())
		})
	}
}
