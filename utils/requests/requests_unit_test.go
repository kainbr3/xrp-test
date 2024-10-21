//go:build unit

package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var unitMethod, unitURL string

func TestCases_Requests_Unit(t *testing.T) {
	unitMethod = "GET"
	unitURL = "https://httpbin.org/get"

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"Success creating a request with given parameters", testBuildRequest},
		{"Success executing a http request", testExecuteRequest},
		{"Failure executing a http request with error", testFailExecuteRequest},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func testBuildRequest(t *testing.T) {
	t.Log("TestBuildRequest - Testing a success clause for creating a request instance with given parameters")
	type payload struct {
		Property int
	}

	data := payload{Property: 1}
	parameters := map[string]interface{}{
		"headers":      map[string]string{"x-api-key": "123"},
		"payload":      data,
		"user-agent":   "Safari/537.3",
		"content-type": "application/json",
	}

	req, err := buildRequest(context.Background(), "GET", "https://httpbin.org/get", parameters)
	// validates request base
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, req.URL.String(), "https://httpbin.org/get")

	// validates request headers
	assert.Equal(t, req.Header["X-Api-Key"][0], string(parameters["headers"].(map[string]string)["x-api-key"]))
	assert.Equal(t, req.Header["User-Agent"][0], parameters["user-agent"])
	assert.Equal(t, req.Header["Content-Type"][0], parameters["content-type"])

	// validates request body
	bodyBytes, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var reqBodyData payload
	err = json.Unmarshal(bodyBytes, &reqBodyData)
	defer req.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, reqBodyData.Property, data.Property)
}

func testExecuteRequest(t *testing.T) {
	t.Log("TestExecuteRequest - Testing a success clause for executing a successful http request")
	req, err := buildRequest(context.Background(), unitMethod, unitURL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)

	ExecuteMock = func() error { return nil }
	err = Execute(context.Background(), unitMethod, unitURL, nil, nil)
	assert.NoError(t, err)
}

func testFailExecuteRequest(t *testing.T) {
	t.Log("TestFailExecuteRequest - Testing a failure clause for executing a http request with error")
	req, err := buildRequest(context.Background(), unitMethod, unitURL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	errorMessage := "error trying execute a http request with statuscode: 400 - endpoint: https://httpbin.org/get - with error: Bad Request"

	ExecuteMock = func() error {
		return fmt.Errorf(errorMessage)
	}

	err = Execute(context.Background(), unitMethod, unitURL, nil, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errorMessage)
}
