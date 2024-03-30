package testhelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"

	"example.com/sample/commonlib/auth"
	"github.com/labstack/echo/v4"
)

type mockRequestKey struct{}

// IsMockContext reports whether this context is from a mock HTTP request or not
func IsMockContext(ctx context.Context) bool {
	return ctx.Value(mockRequestKey{}) != nil
}

// RequestBuilder helps manage the boilerplate of constructing echo contexts while testing REST controllers
type RequestBuilder struct {
	method        string
	path          string
	body          any
	serializeBody bool
	pathParams    map[string]string
	auth          *auth.CustomClaims
}

// NewRequest starts a new request builder with the bare minimum required data
func NewRequest(method string, path string) RequestBuilder {
	return RequestBuilder{
		method: method,
		path:   path,
	}
}

// WithBody adds a JSON body that should be serialized before adding it to the HTTP request being built
func (rb RequestBuilder) WithBody(body any) RequestBuilder {
	rb.body = body
	rb.serializeBody = true
	return rb
}

// WithRawBody adds a raw array of bytes to use for the HTTP request body to the request being built
func (rb RequestBuilder) WithRawBody(body []byte) RequestBuilder {
	rb.body = body
	return rb
}

// WithPathParams registers path parameters and their values with the constructed echo context.
//
// For example, if your route accepts the path /variables/:id, you might provide the map "id: 123" to provide the
// value for "id".
func (rb RequestBuilder) WithPathParams(pathParams map[string]string) RequestBuilder {
	rb.pathParams = pathParams
	return rb
}

// WithAuth adds the specified authentication information to the request being built
func (rb RequestBuilder) WithAuth(authToken auth.CustomClaims) RequestBuilder {
	rb.auth = &authToken
	return rb
}

// Build constructs an echo.Context containing the constructed request and a httptest.ResponseRecorder which
// captures the response a controller passes through the echo.Context
func (rb RequestBuilder) Build() (echo.Context, *httptest.ResponseRecorder, error) {
	var body []byte
	if rb.body != nil {
		if rb.serializeBody {
			var serializeErr error
			body, serializeErr = json.Marshal(rb.body)
			if serializeErr != nil {
				return nil, nil, fmt.Errorf("failed to serialize json body of http request: %w", serializeErr)
			}
		} else {
			body = rb.body.([]byte)
		}
	}

	request := httptest.NewRequest(rb.method, rb.path, bytes.NewBuffer(body))
	// Mark the request's context as a mock context
	request = request.WithContext(context.WithValue(request.Context(), mockRequestKey{}, true))

	if rb.body != nil {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	responseRecorder := httptest.NewRecorder()
	e := echo.New()

	ctx := e.NewContext(request, responseRecorder)

	if rb.auth != nil {
		ctx.Set("user", *rb.auth)
	}
	if rb.pathParams != nil {
		var paramNames []string
		var paramValues []string

		for key, value := range rb.pathParams {
			paramNames = append(paramNames, key)
			paramValues = append(paramValues, value)
		}

		ctx.SetParamNames(paramNames...)
		ctx.SetParamValues(paramValues...)
	}

	return ctx, responseRecorder, nil
}
