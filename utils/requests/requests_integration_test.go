//go:build integration

package requests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCases_Requests_Integration(t *testing.T) {
	unitMethod = "GET"
	unitURL = "https://httpbin.org/get"

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"Success executing a http request", testIntegrationExecuteRequest},
		{"Failure executing a http request with error", testFailIntegrationExecuteRequest},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func testIntegrationExecuteRequest(t *testing.T) {
	t.Log("TestIntegrationExecuteRequest - Testing a success clause for executing a successful http request")

	// creates a runtime struct and assigns to a variable to receive the response from the coingecko ping request
	result := struct {
		Args struct {
		} `json:"args"`
		Headers struct {
			Accept                  string `json:"Accept"`
			AcceptEncoding          string `json:"Accept-Encoding"`
			AcceptLanguage          string `json:"Accept-Language"`
			Host                    string `json:"Host"`
			Priority                string `json:"Priority"`
			SecChUa                 string `json:"Sec-Ch-Ua"`
			SecChUaMobile           string `json:"Sec-Ch-Ua-Mobile"`
			SecChUaPlatform         string `json:"Sec-Ch-Ua-Platform"`
			SecFetchDest            string `json:"Sec-Fetch-Dest"`
			SecFetchMode            string `json:"Sec-Fetch-Mode"`
			SecFetchSite            string `json:"Sec-Fetch-Site"`
			SecFetchUser            string `json:"Sec-Fetch-User"`
			UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
			UserAgent               string `json:"User-Agent"`
			XAmznTraceID            string `json:"X-Amzn-Trace-Id"`
		} `json:"headers"`
		Origin string `json:"origin"`
		URL    string `json:"url"`
	}{}

	err := Execute(context.Background(), unitMethod, unitURL, &result, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.URL, "https://httpbin.org/get")
}

func testFailIntegrationExecuteRequest(t *testing.T) {
	t.Log("TestFailIntegrationExecuteRequest - Testing a failure clause for executing a http request with error")
	err := Execute(context.Background(), "GET", "https://invalid", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), `error trying execute a http request on: https://invalid - with error: Get "https://invalid": dial tcp: lookup invalid: no such host`)
}
