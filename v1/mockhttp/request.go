package mockhttp

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

type Request struct {
	W *httptest.ResponseRecorder
	R *http.Request
}

// NewRequest creates a wrapper object around objects necessary to do a mock http request
func NewRequest(method, path, body string) *Request {
	return &Request{
		W: httptest.NewRecorder(),
		R: httptest.NewRequest(method, path, bytes.NewBuffer([]byte(body))),
	}
}

func (r *Request) WithContext(ctx context.Context) *Request {
	r.R = r.R.WithContext(ctx)
	return r
}

// Inserts each key/value pair into context that gets inserted into the Request
func (r *Request) WithValues(vals map[string]interface{}) *Request {
	ctx := r.R.Context()
	for key, val := range vals {
		ctx = context.WithValue(ctx, key, val)
	}
	r.R = r.R.WithContext(ctx)
	return r
}

// With path params sets path params for the request
func (r *Request) WithPathParams(vals map[string]string) *Request {
	ctx := r.R.Context()
	rctx := chi.NewRouteContext()
	for key, val := range vals {
		rctx.URLParams.Add(key, val)
	}
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	r.R = r.R.WithContext(ctx)
	return r
}

// SetHeader sets HTTP Header for wrapper object
func (r *Request) SetHeader(key, value string) *Request {
	r.R.Header.Set(key, value)
	return r
}

func (r *Request) Result() *http.Response {
	return r.W.Result()
}
