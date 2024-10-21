package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	nurl "net/url"
	"strings"
)

func buildRequest(ctx context.Context, method, url string, parameters map[string]any) (*http.Request, error) {
	// initializes the http request variables
	var req *http.Request
	var body io.Reader
	var err error

	// verify if there is a valid payload, then add it to the body object
	if payload, exists := parameters["payload"]; exists {
		data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("error trying to marshal payload, with error: %s", err.Error())
			body = nil
		}

		body = bytes.NewReader(data)
	}

	// verify if there is a valid form data, then add it to the body object
	if payload, exists := parameters["form"]; exists {
		if payloadMap, ok := payload.(map[string]string); ok {
			formData := nurl.Values{}
			for key, value := range payloadMap {
				formData.Set(key, value)
			}
			body = strings.NewReader(formData.Encode())
		} else {
			log.Printf("payload is not a valid map[string]string")
			body = nil
		}
	}

	// new instance of http.Request
	req, err = http.NewRequestWithContext(ctx, method, url, body)

	//verify if there is any valid header, then add it to the request object
	if headers, exists := parameters["headers"]; exists {
		for key, value := range headers.(map[string]string) {
			// validate if there is a valid header and set it up to the request
			if key != "" && value != "" {
				req.Header.Add(key, value)
			}
		}
	}

	// verify if there is a user agent, then add it to the request object
	if userAgent, exists := parameters["user-agent"]; exists {
		req.Header.Set("User-Agent", userAgent.(string))
	}

	// verify if there is a content-type, then add it to the request object
	if contentType, exists := parameters["content-type"]; exists {
		req.Header.Set("Content-Type", contentType.(string))
	} else { // if none provided, sets the default
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	}

	return req, err
}

var ExecuteMock func() error

// Execute - executes a http request with the given parameters
func Execute(ctx context.Context, method, url string, target any, parameters map[string]any) error {
	// returns the mock result if it is set
	if ExecuteMock != nil {
		return ExecuteMock()
	}

	req, err := buildRequest(ctx, method, url, parameters)
	if err != nil {
		return fmt.Errorf("error while creating the request - with error: %s", err.Error())
	}

	// execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil || !(resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusIMUsed) {
		if err != nil {
			// concacts the known error message with the received url and response error and return it
			return fmt.Errorf(" error: %v calling endpoint: %s", err.Error(), url)
		} else {
			var message any
			json.NewDecoder(resp.Body).Decode(&message)
			// concacts the known error message with the received url and response return it
			return fmt.Errorf("error with status code: %d with message: %v on endpoint: %s", resp.StatusCode, message, url)
		}
	}

	// recommended according to documentation to always close the
	// body close immediately after checking the error
	defer resp.Body.Close()

	// decodes the response body into the received target object
	// when the response does not contain a payload, target sould be nil
	if target != nil {
		json.NewDecoder(resp.Body).Decode(&target)
	}

	// return an empty error
	return nil
}
