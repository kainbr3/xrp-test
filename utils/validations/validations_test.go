//go:build unit

package validations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestCases_Validation_Unit(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"Success validating a valid struct", testValidateValidStruct},
		{"Failure validating an invalid struct", testValidateInvalidStruct},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func testValidateValidStruct(t *testing.T) {
	t.Log("testValidateValidStruct - Testing a success clause for validating a valid struct")
	validStruct := TestStruct{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}
	err := Validate(validStruct)
	assert.NoError(t, err)
}

func testValidateInvalidStruct(t *testing.T) {
	t.Log("testValidateInvalidStruct - Testing a failure clause for validating an invalid struct")
	invalidStruct := TestStruct{
		Name:  "",
		Email: "invalid-email",
	}
	err := Validate(invalidStruct)
	assert.Error(t, err)
}
